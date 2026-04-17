package database

import (
	"fmt"
	"log"
	"time"
	"videohub/backend/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Path             string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  int
}

func NewSQLite(cfg DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&model.Media{},
		&model.Season{},
		&model.Episode{},
		&model.WatchHistory{},
		&model.Favorite{},
		&model.Setting{},
		&model.ScanTask{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func CreateIndexes(db *gorm.DB) error {
	log.Println("Creating database indexes...")

	err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_media_type ON media(type);
		CREATE INDEX IF NOT EXISTS idx_media_year ON media(year);
		CREATE INDEX IF NOT EXISTS idx_media_rating ON media(rating);
		CREATE INDEX IF NOT EXISTS idx_media_status ON media(status);
		CREATE INDEX IF NOT EXISTS idx_media_tmdb_id ON media(tmdb_id);
		CREATE INDEX IF NOT EXISTS idx_episodes_series ON episodes(series_id, season_number);
		CREATE INDEX IF NOT EXISTS idx_history_media ON watch_history(media_id);
		CREATE INDEX IF NOT EXISTS idx_history_episode ON watch_history(episode_id);
		CREATE INDEX IF NOT EXISTS idx_history_watched ON watch_history(last_watched);
		CREATE INDEX IF NOT EXISTS idx_favorite_media ON favorites(media_id);
		CREATE INDEX IF NOT EXISTS idx_scan_task_status ON scan_tasks(status);
		CREATE INDEX IF NOT EXISTS idx_scan_task_folder ON scan_tasks(folder_path);
	`).Error

	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
