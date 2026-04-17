package repository

import (
	"context"
	"videohub/backend/internal/model"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *FavoriteRepository) GetByMedia(ctx context.Context, mediaID int64) (*model.Favorite, error) {
	var favorite model.Favorite
	err := r.db.WithContext(ctx).Where("media_id = ?", mediaID).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *FavoriteRepository) List(ctx context.Context, page, pageSize int) ([]model.Favorite, int64, error) {
	var favoriteList []model.Favorite
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Favorite{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Preload("Media").Order("created_at DESC").Find(&favoriteList).Error
	return favoriteList, total, err
}

func (r *FavoriteRepository) Delete(ctx context.Context, mediaID int64) error {
	return r.db.WithContext(ctx).Where("media_id = ?", mediaID).Delete(&model.Favorite{}).Error
}
