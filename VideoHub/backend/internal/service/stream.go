package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"videohub/backend/internal/model"
	"videohub/backend/internal/repository"
)

var (
	ErrTranscodeFailed = errors.New("transcode failed")
)

type TranscodeOptions struct {
	Quality      string
	VideoCodec   string
	AudioCodec   string
	VideoBitrate string
	AudioBitrate string
	Resolution   string
	HardwareAccel bool
}

var QualityPresets = map[string]TranscodeOptions{
	"4k": {
		VideoCodec:   "h265",
		VideoBitrate: "15000k",
		Resolution:   "3840x2160",
	},
	"1080p": {
		VideoCodec:   "h264",
		VideoBitrate: "5000k",
		Resolution:   "1920x1080",
	},
	"720p": {
		VideoCodec:   "h264",
		VideoBitrate: "2500k",
		Resolution:   "1280x720",
	},
	"480p": {
		VideoCodec:   "h264",
		VideoBitrate: "1000k",
		Resolution:   "854x480",
	},
}

type StreamService struct {
	repo      *repository.Repository
	hlsDir    string
	mutex     sync.Mutex
	transcodeTasks map[string]*TranscodeTask
}

type TranscodeTask struct {
	ID        string
	MediaID   int64
	EpisodeID *int64
	Status    string
	Progress  float64
	Command   *exec.Cmd
}

func NewStreamService(repo *repository.Repository, hlsDir string) *StreamService {
	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create HLS directory: %v\n", err)
	}

	return &StreamService{
		repo:      repo,
		hlsDir:    hlsDir,
		transcodeTasks: make(map[string]*TranscodeTask),
	}
}

func (s *StreamService) GetStreamURL(ctx context.Context, mediaID int64, episodeID *int64, quality string) (string, error) {
	var media *model.Media
	var filePath string

	if episodeID != nil {
		episode, err := s.repo.Episode.GetByID(ctx, *episodeID)
		if err != nil {
			return "", fmt.Errorf("failed to get episode: %w", err)
		}
		filePath = episode.Path
		media, err = s.repo.Media.GetByID(ctx, episode.SeriesID)
		if err != nil {
			return "", fmt.Errorf("failed to get media: %w", err)
		}
	} else {
		var err error
		media, err = s.repo.Media.GetByID(ctx, mediaID)
		if err != nil {
			return "", fmt.Errorf("failed to get media: %w", err)
		}
		filePath = media.Path
	}

	if quality == "" {
		quality = "1080p"
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("media file not found: %s", filePath)
	}

	// 生成HLS流
	taskID, err := s.StartTranscodeTask(ctx, mediaID, episodeID, &TranscodeOptions{
		Quality: quality,
	})
	if err != nil {
		return "", fmt.Errorf("failed to start transcode: %w", err)
	}

	// 构建播放URL
	playURL := fmt.Sprintf("/hls/%s/playlist.m3u8", taskID)
	return playURL, nil
}

func (s *StreamService) Transcode(ctx context.Context, inputPath string, outputPath string, options *TranscodeOptions) error {
	// 获取转码预设
	preset, ok := QualityPresets[options.Quality]
	if !ok {
		preset = QualityPresets["1080p"]
	}

	// 构建FFmpeg命令
	cmdArgs := []string{
		"-i", inputPath,
		"-profile:v", "baseline",
		"-level", "3.0",
		"-s", preset.Resolution,
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		outputPath,
	}

	// 添加硬件加速
	if options.HardwareAccel {
		cmdArgs = append([]string{"-hwaccel", "cuda"}, cmdArgs...)
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}

func (s *StreamService) StartTranscodeTask(ctx context.Context, mediaID int64, episodeID *int64, options *TranscodeOptions) (string, error) {
	taskID := fmt.Sprintf("%d_%d", mediaID, time.Now().Unix())
	if episodeID != nil {
		taskID = fmt.Sprintf("%d_%d_%d", mediaID, *episodeID, time.Now().Unix())
	}

	taskDir := filepath.Join(s.hlsDir, taskID)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create task directory: %w", err)
	}

	// 获取媒体文件路径
	var filePath string
	if episodeID != nil {
		episode, err := s.repo.Episode.GetByID(ctx, *episodeID)
		if err != nil {
			return "", fmt.Errorf("failed to get episode: %w", err)
		}
		filePath = episode.Path
	} else {
		media, err := s.repo.Media.GetByID(ctx, mediaID)
		if err != nil {
			return "", fmt.Errorf("failed to get media: %w", err)
		}
		filePath = media.Path
	}

	// 构建FFmpeg命令
	playlistPath := filepath.Join(taskDir, "playlist.m3u8")
	cmdArgs := []string{
		"-i", filePath,
		"-profile:v", "baseline",
		"-level", "3.0",
		"-s", QualityPresets[options.Quality].Resolution,
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		playlistPath,
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	task := &TranscodeTask{
		ID:        taskID,
		MediaID:   mediaID,
		EpisodeID: episodeID,
		Status:    "running",
		Progress:  0,
		Command:   cmd,
	}

	s.mutex.Lock()
	s.transcodeTasks[taskID] = task
	s.mutex.Unlock()

	// 启动后台进程监控转码进度
	go func() {
		cmd.Wait()
		s.mutex.Lock()
		if task, exists := s.transcodeTasks[taskID]; exists {
			task.Status = "completed"
			task.Progress = 100
		}
		s.mutex.Unlock()
	}()

	return taskID, nil
}

func (s *StreamService) GetTranscodeProgress(ctx context.Context, taskID string) (float64, error) {
	s.mutex.Lock()
	task, exists := s.transcodeTasks[taskID]
	s.mutex.Unlock()

	if !exists {
		return 0, errors.New("transcode task not found")
	}

	return task.Progress, nil
}

func (s *StreamService) CancelTranscodeTask(ctx context.Context, taskID string) error {
	s.mutex.Lock()
	task, exists := s.transcodeTasks[taskID]
	s.mutex.Unlock()

	if !exists {
		return errors.New("transcode task not found")
	}

	if task.Command != nil && task.Command.Process != nil {
		if err := task.Command.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill transcode process: %w", err)
		}
	}

	s.mutex.Lock()
	delete(s.transcodeTasks, taskID)
	s.mutex.Unlock()

	// 清理临时文件
	taskDir := filepath.Join(s.hlsDir, taskID)
	os.RemoveAll(taskDir)

	return nil
}

func (s *StreamService) CleanupOldStreams(ctx context.Context) error {
	// 清理超过24小时的转码任务
	threshold := time.Now().Add(-24 * time.Hour)

	s.mutex.Lock()
	for taskID, task := range s.transcodeTasks {
		if task.Status == "completed" && time.Since(threshold) > 0 {
			delete(s.transcodeTasks, taskID)
			// 清理文件
			taskDir := filepath.Join(s.hlsDir, taskID)
			os.RemoveAll(taskDir)
		}
	}
	s.mutex.Unlock()

	return nil
}
