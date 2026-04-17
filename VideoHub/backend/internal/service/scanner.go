package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"videohub/backend/internal/model"
	"videohub/backend/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrUnknownMediaType = errors.New("unknown media type")
	ErrScanFailed       = errors.New("scan failed")
)

type MediaFile struct {
	Path     string
	Type     model.MediaType
	Title    string
	Year     string
	Quality  string
	Is3D     bool
	Size     int64
	Modified time.Time
}

type ScanTask struct {
	ID             int64
	FolderPath     string
	Status         string
	TotalFiles     int
	ProcessedFiles int
	Progress       float64
	ErrorMessage   string
	StartedAt      time.Time
	CompletedAt    time.Time
}

type ScannerService struct {
	repo      *repository.Repository
	db        *gorm.DB
	videoExts map[string]bool
	moviePatterns []*regexp.Regexp
	tvPatterns    []*regexp.Regexp
	mutex     sync.Mutex
	tasks     map[int64]*ScanTask
}

func NewScannerService(repo *repository.Repository, db *gorm.DB) *ScannerService {
	return &ScannerService{
		repo:      repo,
		db:        db,
		videoExts: getVideoExtensions(),
		moviePatterns: getMoviePatterns(),
		tvPatterns:    getTVPatterns(),
		tasks:     make(map[int64]*ScanTask),
	}
}

func getVideoExtensions() map[string]bool {
	return map[string]bool{
		".mkv": true, ".mp4": true, ".avi": true, ".mov": true,
		".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
		".mpg": true, ".mpeg": true, ".3gp": true,
	}
}

func getMoviePatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)(.+?)[.\s]+(19|20)\d{2}[.\s].*\.(mkv|mp4|avi)`),
		regexp.MustCompile(`(?i)(.+?)[.\s]+(2160p|1080p|720p|480p)[.\s].*\.(mkv|mp4|avi)`),
	}
}

func getTVPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)S(\d{1,2})E(\d{1,2})`),
		regexp.MustCompile(`(?i)Season[.\s]*(\d{1,2})[.\s]*Episode[.\s]*(\d{1,2})`),
		regexp.MustCompile(`(?i)(\d{1,2})x(\d{1,2})`),
	}
}

func (s *ScannerService) StartScan(ctx context.Context, folderPath string, force bool) (*ScanTask, error) {
	task := &model.ScanTask{
		FolderPath: folderPath,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("failed to create scan task: %w", err)
	}

	scanTask := &ScanTask{
		ID:         task.ID,
		FolderPath: folderPath,
		Status:     "scanning",
		StartedAt:  time.Now(),
	}

	s.mutex.Lock()
	s.tasks[task.ID] = scanTask
	s.mutex.Unlock()

	go s.scanFolderAsync(ctx, scanTask, force)

	return scanTask, nil
}

func (s *ScannerService) scanFolderAsync(ctx context.Context, task *ScanTask, force bool) {
	defer func() {
		s.mutex.Lock()
		delete(s.tasks, task.ID)
		s.mutex.Unlock()
	}()

	files, err := s.ScanFolder(ctx, task.FolderPath)
	if err != nil {
		task.Status = "failed"
		task.ErrorMessage = err.Error()
		s.updateTaskStatus(ctx, task)
		return
	}

	task.TotalFiles = len(files)
	task.ProcessedFiles = 0
	task.Progress = 0

	for i, file := range files {
		if err := s.processMediaFile(ctx, file, force); err != nil {
			// 记录错误但继续处理其他文件
			fmt.Printf("Failed to process %s: %v\n", file.Path, err)
		}

		task.ProcessedFiles = i + 1
		task.Progress = float64(task.ProcessedFiles) / float64(task.TotalFiles) * 100
		s.updateTaskStatus(ctx, task)

		// 检查上下文是否取消
		if ctx.Err() != nil {
			task.Status = "cancelled"
			s.updateTaskStatus(ctx, task)
			return
		}
	}

	task.Status = "completed"
	task.CompletedAt = time.Now()
	s.updateTaskStatus(ctx, task)
}

func (s *ScannerService) updateTaskStatus(ctx context.Context, task *ScanTask) {
	dbTask := &model.ScanTask{
		ID:             task.ID,
		Status:         task.Status,
		TotalFiles:     task.TotalFiles,
		ProcessedFiles: task.ProcessedFiles,
		ErrorMessage:   task.ErrorMessage,
		StartedAt:      task.StartedAt,
		CompletedAt:    task.CompletedAt,
	}
	s.db.Save(dbTask)
}

func (s *ScannerService) GetScanProgress(ctx context.Context, taskID int64) (*ScanTask, error) {
	s.mutex.Lock()
	task, exists := s.tasks[taskID]
	s.mutex.Unlock()

	if exists {
		return task, nil
	}

	var dbTask model.ScanTask
	if err := s.db.First(&dbTask, taskID).Error; err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	return &ScanTask{
		ID:             dbTask.ID,
		FolderPath:     dbTask.FolderPath,
		Status:         dbTask.Status,
		TotalFiles:     dbTask.TotalFiles,
		ProcessedFiles: dbTask.ProcessedFiles,
		Progress:       float64(dbTask.ProcessedFiles) / float64(dbTask.TotalFiles) * 100,
		ErrorMessage:   dbTask.ErrorMessage,
		StartedAt:      dbTask.StartedAt,
		CompletedAt:    dbTask.CompletedAt,
	}, nil
}

func (s *ScannerService) CancelScan(ctx context.Context, taskID int64) error {
	s.mutex.Lock()
	task, exists := s.tasks[taskID]
	s.mutex.Unlock()

	if exists {
		task.Status = "cancelled"
		s.updateTaskStatus(ctx, task)
	}

	return nil
}

func (s *ScannerService) ScanFolder(ctx context.Context, folderPath string) ([]MediaFile, error) {
	var files []MediaFile

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 跳过隐藏目录
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if s.IsMediaFile(path) {
			mediaType, err := s.GetMediaType(path)
			if err != nil {
				// 跳过无法识别类型的文件
				return nil
			}

			title, year, quality, is3D := s.ParseMediaInfo(info.Name())
			files = append(files, MediaFile{
				Path:     path,
				Type:     mediaType,
				Title:    title,
				Year:     year,
				Quality:  quality,
				Is3D:     is3D,
				Size:     info.Size(),
				Modified: info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan folder: %w", err)
	}

	return files, nil
}

func (s *ScannerService) IsMediaFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return s.videoExts[ext]
}

func (s *ScannerService) GetMediaType(filename string) (model.MediaType, error) {
	basename := filepath.Base(filename)

	for _, pattern := range s.tvPatterns {
		if pattern.MatchString(basename) {
			return model.MediaTypeTV, nil
		}
	}

	for _, pattern := range s.moviePatterns {
		if pattern.MatchString(basename) {
			return model.MediaTypeMovie, nil
		}
	}

	return "", ErrUnknownMediaType
}

func (s *ScannerService) ParseMediaInfo(filename string) (title, year, quality string, is3D bool) {
	re := regexp.MustCompile(`(?i)(.+?)[.\s]+(19|20)\d{2}[.\s]`)
	if matches := re.FindStringSubmatch(filename); len(matches) > 2 {
		title = strings.ReplaceAll(matches[1], ".", " ")
		year = matches[2]
	}

	qualityPattern := regexp.MustCompile(`(?i)(2160p|1080p|720p|480p|4k|hd)`)
	if match := qualityPattern.FindString(filename); match != "" {
		quality = strings.ToUpper(match)
	}

	is3D = regexp.MustCompile(`(?i)3d|half-ou|half-su`).MatchString(filename)
	return
}

func (s *ScannerService) processMediaFile(ctx context.Context, file MediaFile, force bool) error {
	// 检查文件是否已存在
	existing, err := s.repo.Media.GetByPath(ctx, file.Path)
	if err == nil && !force {
		// 文件已存在且不需要强制更新
		return nil
	}

	if existing != nil {
		// 更新现有媒体
		existing.Title = file.Title
		existing.Year = 0
		if file.Year != "" {
			fmt.Sscanf(file.Year, "%d", &existing.Year)
		}
		existing.Type = file.Type
		existing.Status = model.MediaStatusPending
		return s.repo.Media.Update(ctx, existing)
	}

	// 创建新媒体
	year := 0
	if file.Year != "" {
		fmt.Sscanf(file.Year, "%d", &year)
	}

	media := &model.Media{
		Type:   file.Type,
		Title:  file.Title,
		Year:   year,
		Path:   file.Path,
		Status: model.MediaStatusPending,
	}

	return s.repo.Media.Create(ctx, media)
}
