package repository

import (
	"context"
	"videohub/backend/internal/model"

	"gorm.io/gorm"
)

type EpisodeRepository struct {
	db *gorm.DB
}

func NewEpisodeRepository(db *gorm.DB) *EpisodeRepository {
	return &EpisodeRepository{db: db}
}

func (r *EpisodeRepository) Create(ctx context.Context, episode *model.Episode) error {
	return r.db.WithContext(ctx).Create(episode).Error
}

func (r *EpisodeRepository) GetByID(ctx context.Context, id int64) (*model.Episode, error) {
	var episode model.Episode
	err := r.db.WithContext(ctx).First(&episode, id).Error
	if err != nil {
		return nil, err
	}
	return &episode, nil
}

func (r *EpisodeRepository) GetBySeries(ctx context.Context, seriesID int64, season int) ([]model.Episode, error) {
	var episodes []model.Episode
	query := r.db.WithContext(ctx).Where("series_id = ?", seriesID)
	if season > 0 {
		query = query.Where("season_number = ?", season)
	}
	err := query.Order("season_number, episode_number").Find(&episodes).Error
	return episodes, err
}

func (r *EpisodeRepository) Update(ctx context.Context, episode *model.Episode) error {
	return r.db.WithContext(ctx).Save(episode).Error
}

func (r *EpisodeRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Episode{}, id).Error
}
