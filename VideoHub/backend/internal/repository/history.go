package repository

import (
	"context"
	"videohub/backend/internal/model"

	"gorm.io/gorm"
)

type HistoryRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) Create(ctx context.Context, history *model.WatchHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *HistoryRepository) GetByMedia(ctx context.Context, mediaID int64, episodeID *int64) (*model.WatchHistory, error) {
	var history model.WatchHistory
	query := r.db.WithContext(ctx).Where("media_id = ?", mediaID)
	if episodeID != nil {
		query = query.Where("episode_id = ?", *episodeID)
	} else {
		query = query.Where("episode_id IS NULL")
	}
	err := query.First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *HistoryRepository) List(ctx context.Context, page, pageSize int) ([]model.WatchHistory, int64, error) {
	var historyList []model.WatchHistory
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.WatchHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Preload("Media").Order("last_watched DESC").Find(&historyList).Error
	return historyList, total, err
}

func (r *HistoryRepository) Update(ctx context.Context, history *model.WatchHistory) error {
	return r.db.WithContext(ctx).Save(history).Error
}

func (r *HistoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.WatchHistory{}, id).Error
}

func (r *HistoryRepository) ClearAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Delete(&model.WatchHistory{}).Error
}
