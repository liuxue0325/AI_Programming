package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"videohub/backend/internal/model"
	"videohub/backend/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrScrapeFailed     = errors.New("scrape failed")
	ErrExternalAPIError = errors.New("external API error")
)

type TMDBSearchResult struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	OriginalTitle string `json:"original_title"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
	Overview    string `json:"overview"`
	VoteAverage float64 `json:"vote_average"`
}

type TMDBSearchResponse struct {
	Results []TMDBSearchResult `json:"results"`
	TotalResults int64 `json:"total_results"`
}

type TMDBDetails struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	OriginalTitle string `json:"original_title"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
	Overview    string `json:"overview"`
	VoteAverage float64 `json:"vote_average"`
	Runtime     int    `json:"runtime"`
	Genres      []struct {
		Name string `json:"name"`
	} `json:"genres"`
	ImdbID      string `json:"imdb_id"`
}

type TMDBClient struct {
	APIKey    string
	BaseURL   string
	Language  string
	HTTPClient *http.Client
}

type ScraperService struct {
	repo   *repository.Repository
	db     *gorm.DB
	tmdb   *TMDBClient
	cache  map[string]interface{}
}

func NewScraperService(repo *repository.Repository, db *gorm.DB, apiKey string) *ScraperService {
	return &ScraperService{
		repo:  repo,
		db:    db,
		tmdb: &TMDBClient{
			APIKey:    apiKey,
			BaseURL:   "https://api.themoviedb.org/3",
			Language:  "zh-CN",
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
		cache: make(map[string]interface{}),
	}
}

func (s *ScraperService) ScrapeMedia(ctx context.Context, mediaID int64, source string) error {
	media, err := s.repo.Media.GetByID(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("failed to get media: %w", err)
	}

	if source == "tmdb" {
		return s.scrapeFromTMDB(ctx, media)
	}

	return fmt.Errorf("unsupported scraper source: %s", source)
}

func (s *ScraperService) ScrapeBatch(ctx context.Context, mediaIDs []int64, source string) error {
	for _, mediaID := range mediaIDs {
		if err := s.ScrapeMedia(ctx, mediaID, source); err != nil {
			// 记录错误但继续处理其他媒体
			fmt.Printf("Failed to scrape media %d: %v\n", mediaID, err)
		}
	}
	return nil
}

func (s *ScraperService) ScrapeAll(ctx context.Context, force bool) error {
	var mediaList []model.Media
	query := s.db.Where("status = ?", model.MediaStatusPending)
	if force {
		query = s.db.Where("status IN (?, ?)", model.MediaStatusPending, model.MediaStatusScraped)
	}

	if err := query.Find(&mediaList).Error; err != nil {
		return fmt.Errorf("failed to get media list: %w", err)
	}

	for _, media := range mediaList {
		if err := s.scrapeFromTMDB(ctx, &media); err != nil {
			// 记录错误但继续处理其他媒体
			fmt.Printf("Failed to scrape media %d: %v\n", media.ID, err)
		}
	}

	return nil
}

func (s *ScraperService) scrapeFromTMDB(ctx context.Context, media *model.Media) error {
	// 搜索媒体
	results, err := s.SearchTMDB(ctx, media.Title, media.Type)
	if err != nil {
		media.Status = model.MediaStatusFailed
		s.repo.Media.Update(ctx, media)
		return fmt.Errorf("failed to search TMDB: %w", err)
	}

	if len(results) == 0 {
		media.Status = model.MediaStatusFailed
		s.repo.Media.Update(ctx, media)
		return errors.New("no results found on TMDB")
	}

	// 获取详情
	bestMatch := results[0]
	details, err := s.GetTMDBDetails(ctx, bestMatch.ID, media.Type)
	if err != nil {
		media.Status = model.MediaStatusFailed
		s.repo.Media.Update(ctx, media)
		return fmt.Errorf("failed to get TMDB details: %w", err)
	}

	// 更新媒体信息
	media.OriginalTitle = details.OriginalTitle
	media.Overview = details.Overview
	media.Rating = details.VoteAverage
	media.Runtime = details.Runtime
	media.TmdbID = details.ID
	media.ImdbID = details.ImdbID

	// 提取年份
	if details.ReleaseDate != "" {
		date, err := time.Parse("2006-01-02", details.ReleaseDate)
		if err == nil {
			media.Year = date.Year()
		}
	}

	// 提取 genres
	var genres []string
	for _, genre := range details.Genres {
		genres = append(genres, genre.Name)
	}
	if len(genres) > 0 {
		media.Genres = fmt.Sprintf("%v", genres)
	}

	// 下载海报和背景图
	if details.PosterPath != "" {
		posterPath, err := s.DownloadPoster(ctx, details.PosterPath, media.ID)
		if err == nil {
			media.PosterPath = posterPath
		}
	}

	if details.BackdropPath != "" {
		backdropPath, err := s.DownloadBackdrop(ctx, details.BackdropPath, media.ID)
		if err == nil {
			media.BackdropPath = backdropPath
		}
	}

	media.Status = model.MediaStatusScraped
	return s.repo.Media.Update(ctx, media)
}

func (s *ScraperService) SearchTMDB(ctx context.Context, query string, mediaType model.MediaType) ([]TMDBSearchResult, error) {
	url := fmt.Sprintf("%s/search/%s?api_key=%s&query=%s&language=%s",
		s.tmdb.BaseURL, string(mediaType), s.tmdb.APIKey, query, s.tmdb.Language)

	resp, err := s.tmdb.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request TMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result TMDBSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Results, nil
}

func (s *ScraperService) GetTMDBDetails(ctx context.Context, tmdbID int64, mediaType model.MediaType) (*TMDBDetails, error) {
	url := fmt.Sprintf("%s/%s/%d?api_key=%s&language=%s&append_to_response=credits",
		s.tmdb.BaseURL, string(mediaType), tmdbID, s.tmdb.APIKey, s.tmdb.Language)

	resp, err := s.tmdb.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request TMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var details TMDBDetails
	if err := json.Unmarshal(body, &details); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &details, nil
}

func (s *ScraperService) DownloadPoster(ctx context.Context, posterPath string, mediaID int64) (string, error) {
	return s.downloadImage(ctx, posterPath, mediaID, "posters")
}

func (s *ScraperService) DownloadBackdrop(ctx context.Context, backdropPath string, mediaID int64) (string, error) {
	return s.downloadImage(ctx, backdropPath, mediaID, "backdrops")
}

func (s *ScraperService) downloadImage(ctx context.Context, imagePath string, mediaID int64, folder string) (string, error) {
	baseURL := "https://image.tmdb.org/t/p/w500"
	url := baseURL + imagePath

	dir := filepath.Join("./data", folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filename := fmt.Sprintf("%d%s", mediaID, filepath.Ext(imagePath))
	path := filepath.Join(dir, filename)

	resp, err := s.tmdb.HTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return path, nil
}
