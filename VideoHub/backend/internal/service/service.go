package service

import (
	"videohub/backend/internal/repository"

	"gorm.io/gorm"
)

type Service struct {
	Media    *MediaService
	Scanner  *ScannerService
	Scraper  *ScraperService
	Stream   *StreamService
	Setting  *SettingService
}

func NewService(repo *repository.Repository, db *gorm.DB, tmdbAPIKey string, hlsDir string) *Service {
	return &Service{
		Media:    NewMediaService(repo),
		Scanner:  NewScannerService(repo, db),
		Scraper:  NewScraperService(repo, db, tmdbAPIKey),
		Stream:   NewStreamService(repo, hlsDir),
		Setting:  NewSettingService(repo),
	}
}
