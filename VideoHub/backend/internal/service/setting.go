package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"videohub/backend/internal/model"
	"videohub/backend/internal/repository"
)

var (
	ErrSettingNotFound = errors.New("setting not found")
)

type SystemSettings struct {
	Library LibrarySettings `json:"library"`
	Scraper ScraperSettings `json:"scraper"`
	Stream  StreamSettings   `json:"stream"`
	System  SystemSettings   `json:"system"`
}

type LibrarySettings struct {
	MediaFolders     []string `json:"media_folders"`
	ScanOnStartup    bool     `json:"scan_on_startup"`
	AutoScrape       bool     `json:"auto_scrape"`
	ExcludePatterns  []string `json:"exclude_patterns"`
	FileSizeLimit    int64    `json:"file_size_limit"`
}

type ScraperSettings struct {
	TMDBAPIKey          string   `json:"tmdb_api_key"`
	PreferredLanguage  string   `json:"preferred_language"`
	AutoDownloadPoster  bool     `json:"auto_download_poster"`
	AutoDownloadBackdrop bool    `json:"auto_download_backdrop"`
	AutoDownloadSubtitle bool    `json:"auto_download_subtitle"`
	AdditionalSources   []string `json:"additional_sources"`
}

type StreamSettings struct {
	DefaultQuality      string `json:"default_quality"`
	TranscodeEnabled    bool   `json:"transcode_enabled"`
	HardwareAccel       bool   `json:"hardware_acceleration"`
	AdvertisementEnabled bool  `json:"advertisement_enabled"`
	MaxConcurrent       int    `json:"max_concurrent_streams"`
}

type SystemSettings struct {
	Theme        string `json:"theme"`
	Language     string `json:"language"`
	LogLevel     string `json:"log_level"`
	EnableUPNP   bool   `json:"enable_upnp"`
	EnableSharing bool  `json:"enable_sharing"`
}

type SettingService struct {
	repo *repository.Repository
}

func NewSettingService(repo *repository.Repository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetSettings(ctx context.Context) (*SystemSettings, error) {
	settings, err := s.repo.Setting.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// 构建设置映射
	settingMap := make(map[string]string)
	for _, setting := range settings {
		settingMap[setting.Key] = setting.Value
	}

	// 解析设置
	systemSettings := &SystemSettings{
		Library: LibrarySettings{
			MediaFolders:     s.parseStringArray(settingMap, "library.media_folders"),
			ScanOnStartup:    s.parseBool(settingMap, "library.scan_on_startup", true),
			AutoScrape:       s.parseBool(settingMap, "library.auto_scrape", true),
			ExcludePatterns:  s.parseStringArray(settingMap, "library.exclude_patterns"),
			FileSizeLimit:    s.parseInt64(settingMap, "library.file_size_limit", 10737418240),
		},
		Scraper: ScraperSettings{
			TMDBAPIKey:          settingMap["scraper.tmdb_api_key"],
			PreferredLanguage:  settingMap["scraper.preferred_language"],
			AutoDownloadPoster:  s.parseBool(settingMap, "scraper.auto_download_poster", true),
			AutoDownloadBackdrop: s.parseBool(settingMap, "scraper.auto_download_backdrop", true),
			AutoDownloadSubtitle: s.parseBool(settingMap, "scraper.auto_download_subtitle", false),
			AdditionalSources:   s.parseStringArray(settingMap, "scraper.additional_sources"),
		},
		Stream: StreamSettings{
			DefaultQuality:      settingMap["stream.default_quality"],
			TranscodeEnabled:    s.parseBool(settingMap, "stream.transcode_enabled", true),
			HardwareAccel:       s.parseBool(settingMap, "stream.hardware_acceleration", false),
			AdvertisementEnabled: s.parseBool(settingMap, "stream.advertisement_enabled", false),
			MaxConcurrent:       s.parseInt(settingMap, "stream.max_concurrent_streams", 3),
		},
		System: SystemSettings{
			Theme:        settingMap["system.theme"],
			Language:     settingMap["system.language"],
			LogLevel:     settingMap["system.log_level"],
			EnableUPNP:   s.parseBool(settingMap, "system.enable_upnp", false),
			EnableSharing: s.parseBool(settingMap, "system.enable_sharing", false),
		},
	}

	// 设置默认值
	if systemSettings.Scraper.PreferredLanguage == "" {
		systemSettings.Scraper.PreferredLanguage = "zh-CN"
	}
	if systemSettings.Stream.DefaultQuality == "" {
		systemSettings.Stream.DefaultQuality = "1080p"
	}
	if systemSettings.System.Theme == "" {
		systemSettings.System.Theme = "dark"
	}
	if systemSettings.System.Language == "" {
		systemSettings.System.Language = "zh-CN"
	}
	if systemSettings.System.LogLevel == "" {
		systemSettings.System.LogLevel = "info"
	}

	return systemSettings, nil
}

func (s *SettingService) GetSetting(ctx context.Context, key string) (string, error) {
	setting, err := s.repo.Setting.Get(ctx, key)
	if err != nil {
		return "", ErrSettingNotFound
	}
	return setting.Value, nil
}

func (s *SettingService) UpdateSettings(ctx context.Context, settings map[string]interface{}) error {
	for key, value := range settings {
		if err := s.UpdateSetting(ctx, key, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingService) UpdateSetting(ctx context.Context, key string, value interface{}) error {
	var stringValue string
	switch v := value.(type) {
	case string:
		stringValue = v
	case bool:
		stringValue = fmt.Sprintf("%v", v)
	case int:
		stringValue = fmt.Sprintf("%d", v)
	case []string:
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal array: %w", err)
		}
		stringValue = string(data)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		stringValue = string(data)
	}

	description := s.getSettingDescription(key)
	return s.repo.Setting.Set(ctx, key, stringValue, description)
}

func (s *SettingService) ResetSettings(ctx context.Context) error {
	// 获取所有设置
	settings, err := s.repo.Setting.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 删除所有设置
	for _, setting := range settings {
		if err := s.repo.Setting.Delete(ctx, setting.Key); err != nil {
			return fmt.Errorf("failed to delete setting: %w", err)
		}
	}

	// 设置默认值
	defaultSettings := map[string]interface{}{
		"library.scan_on_startup": true,
		"library.auto_scrape": true,
		"scraper.preferred_language": "zh-CN",
		"scraper.auto_download_poster": true,
		"scraper.auto_download_backdrop": true,
		"stream.default_quality": "1080p",
		"stream.transcode_enabled": true,
		"system.theme": "dark",
		"system.language": "zh-CN",
		"system.log_level": "info",
	}

	return s.UpdateSettings(ctx, defaultSettings)
}

func (s *SettingService) ValidateSettings(ctx context.Context, settings *SystemSettings) error {
	// 验证TMDB API Key
	if settings.Scraper.TMDBAPIKey == "" {
		return errors.New("TMDB API key is required")
	}

	// 验证媒体文件夹
	if len(settings.Library.MediaFolders) == 0 {
		return errors.New("at least one media folder is required")
	}

	return nil
}

func (s *SettingService) parseBool(settings map[string]string, key string, defaultValue bool) bool {
	if value, exists := settings[key]; exists {
		if value == "true" {
			return true
		} else if value == "false" {
			return false
		}
	}
	return defaultValue
}

func (s *SettingService) parseInt(settings map[string]string, key string, defaultValue int) int {
	if value, exists := settings[key]; exists {
		if intValue, err := fmt.Sscanf(value, "%d", &intValue); err == nil && intValue > 0 {
			return intValue
		}
	}
	return defaultValue
}

func (s *SettingService) parseInt64(settings map[string]string, key string, defaultValue int64) int64 {
	if value, exists := settings[key]; exists {
		if intValue, err := fmt.Sscanf(value, "%d", &intValue); err == nil && intValue > 0 {
			return int64(intValue)
		}
	}
	return defaultValue
}

func (s *SettingService) parseStringArray(settings map[string]string, key string) []string {
	if value, exists := settings[key]; exists {
		var array []string
		if err := json.Unmarshal([]byte(value), &array); err == nil {
			return array
		}
	}
	return []string{}
}

func (s *SettingService) getSettingDescription(key string) string {
	descriptions := map[string]string{
		"library.media_folders": "Media folders to scan",
		"library.scan_on_startup": "Scan media on startup",
		"library.auto_scrape": "Auto scrape metadata",
		"scraper.tmdb_api_key": "TMDB API key",
		"scraper.preferred_language": "Preferred language",
		"stream.default_quality": "Default streaming quality",
		"stream.transcode_enabled": "Enable transcoding",
		"system.theme": "UI theme",
		"system.language": "System language",
		"system.log_level": "Log level",
	}
	return descriptions[key]
}
