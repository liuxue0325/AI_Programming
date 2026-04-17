package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"videohub/backend/internal/model"
	"videohub/backend/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrMediaNotFound  = errors.New("media not found")
	ErrInvalidParameter = errors.New("invalid parameter")
)

type MediaListRequest struct {
	Type     *model.MediaType
	Year     *int
	Genre    *string
	Status   *model.MediaStatus
	Keyword  *string
	Page     int
	PageSize int
	Sort     string
	Order    string
}

type MediaListResponse struct {
	Items      []model.Media `json:"items"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

type PlayInfo struct {
	MediaID      int64      `json:"media_id"`
	Type         string     `json:"type"`
	PlayURL      string     `json:"play_url"`
	Subtitles    []Subtitle `json:"subtitles"`
	TranscodeNeeded bool     `json:"transcode_needed"`
	Quality      string     `json:"quality"`
	Bitrate      int        `json:"bitrate"`
}

type Subtitle struct {
	ID       int64  `json:"id"`
	Language string `json:"language"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Format   string `json:"format"`
	Path     string `json:"path"`
}

type MediaService struct {
	repo *repository.Repository
}

func NewMediaService(repo *repository.Repository) *MediaService {
	return &MediaService{repo: repo}
}

func (s *MediaService) GetMediaList(ctx context.Context, req *MediaListRequest) (*MediaListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	mediaList, total, err := s.repo.Media.List(ctx, req.Type, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get media list: %w", err)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &MediaListResponse{
		Items:      mediaList,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *MediaService) GetMediaByID(ctx context.Context, id int64) (*model.Media, error) {
	media, err := s.repo.Media.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMediaNotFound
		}
		return nil, fmt.Errorf("failed to get media: %w", err)
	}
	return media, nil
}

func (s *MediaService) GetMediaByPath(ctx context.Context, path string) (*model.Media, error) {
	media, err := s.repo.Media.GetByPath(ctx, path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMediaNotFound
		}
		return nil, fmt.Errorf("failed to get media: %w", err)
	}
	return media, nil
}

func (s *MediaService) DeleteMedia(ctx context.Context, id int64) error {
	_, err := s.GetMediaByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Media.Delete(ctx, id)
}

func (s *MediaService) GetEpisodes(ctx context.Context, seriesID int64, season int) ([]model.Episode, error) {
	return s.repo.Episode.GetBySeries(ctx, seriesID, season)
}

func (s *MediaService) GetEpisodeByID(ctx context.Context, id int64) (*model.Episode, error) {
	episode, err := s.repo.Episode.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("episode not found")
		}
		return nil, fmt.Errorf("failed to get episode: %w", err)
	}
	return episode, nil
}

func (s *MediaService) GetPlayInfo(ctx context.Context, mediaID int64, episodeID *int64) (*PlayInfo, error) {
	media, err := s.GetMediaByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	var filePath string
	if episodeID != nil {
		episode, err := s.GetEpisodeByID(ctx, *episodeID)
		if err != nil {
			return nil, err
		}
		filePath = episode.Path
	} else {
		filePath = media.Path
	}

	subtitles, err := s.GetSubtitles(ctx, mediaID)
	if err != nil {
		subtitles = []Subtitle{}
	}

	return &PlayInfo{
		MediaID:      mediaID,
		Type:         string(media.Type),
		PlayURL:      fmt.Sprintf("/stream/%d", mediaID),
		Subtitles:    subtitles,
		TranscodeNeeded: true, // 简化处理，实际应根据文件格式判断
		Quality:      "1080p",
		Bitrate:      5000000,
	}, nil
}

func (s *MediaService) GetSubtitles(ctx context.Context, mediaID int64) ([]Subtitle, error) {
	media, err := s.GetMediaByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	// 简化实现，实际应扫描字幕文件
	subtitles := []Subtitle{}
	basePath := filepath.Dir(media.Path)
	title := strings.TrimSuffix(filepath.Base(media.Path), filepath.Ext(media.Path))

	// 模拟字幕文件
	subExts := []string{"srt", "ass", "vtt"}
	langs := []string{"zh-CN", "en"}

	for i, lang := range langs {
		for _, ext := range subExts {
			subPath := filepath.Join(basePath, fmt.Sprintf("%s.%s.%s", title, lang, ext))
			subtitles = append(subtitles, Subtitle{
				ID:       int64(i),
				Language: lang,
				Name:     map[string]string{"zh-CN": "简体中文", "en": "English"}[lang],
				URL:      fmt.Sprintf("/api/media/%d/subtitles/%d", mediaID, i),
				Format:   ext,
				Path:     subPath,
			})
		}
	}

	return subtitles, nil
}

func (s *MediaService) SearchMedia(ctx context.Context, keyword string) ([]model.Media, error) {
	if keyword == "" {
		return nil, ErrInvalidParameter
	}

	return s.repo.Media.Search(ctx, keyword)
}

func (s *MediaService) GetWatchHistory(ctx context.Context, page, pageSize int) ([]model.WatchHistory, int64, error) {
	return s.repo.History.List(ctx, page, pageSize)
}

func (s *MediaService) AddWatchHistory(ctx context.Context, mediaID int64, episodeID *int64, progress float64, completed bool) error {
	history, err := s.repo.History.GetByMedia(ctx, mediaID, episodeID)
	if err != nil {
		// 创建新历史记录
		history = &model.WatchHistory{
			MediaID:     mediaID,
			EpisodeID:   episodeID,
			Progress:    progress,
			Completed:   completed,
			LastWatched: time.Now(),
		}
		return s.repo.History.Create(ctx, history)
	}

	// 更新现有历史记录
	history.Progress = progress
	history.Completed = completed
	history.LastWatched = time.Now()
	return s.repo.History.Update(ctx, history)
}

func (s *MediaService) ClearWatchHistory(ctx context.Context) error {
	return s.repo.History.ClearAll(ctx)
}

func (s *MediaService) GetFavorites(ctx context.Context, page, pageSize int) ([]model.Favorite, int64, error) {
	return s.repo.Favorite.List(ctx, page, pageSize)
}

func (s *MediaService) AddFavorite(ctx context.Context, mediaID int64) error {
	// 检查媒体是否存在
	_, err := s.GetMediaByID(ctx, mediaID)
	if err != nil {
		return err
	}

	// 检查是否已收藏
	_, err = s.repo.Favorite.GetByMedia(ctx, mediaID)
	if err == nil {
		// 已收藏，无需重复添加
		return nil
	}

	// 创建收藏
	favorite := &model.Favorite{
		MediaID: mediaID,
	}
	return s.repo.Favorite.Create(ctx, favorite)
}

func (s *MediaService) DeleteFavorite(ctx context.Context, mediaID int64) error {
	return s.repo.Favorite.Delete(ctx, mediaID)
}
