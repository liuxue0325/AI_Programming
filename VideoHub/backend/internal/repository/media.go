package repository

import (
	"context"
	"videohub/backend/internal/model"

	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(ctx context.Context, media *model.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *MediaRepository) GetByID(ctx context.Context, id int64) (*model.Media, error) {
	var media model.Media
	err := r.db.WithContext(ctx).Preload("Seasons").Preload("Seasons.Episodes").First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) GetByPath(ctx context.Context, path string) (*model.Media, error) {
	var media model.Media
	err := r.db.WithContext(ctx).Where("path = ?", path).First(&media).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) List(ctx context.Context, mediaType *model.MediaType, page, pageSize int) ([]model.Media, int64, error) {
	var mediaList []model.Media
	var total int64

	query := r.db.WithContext(ctx)
	if mediaType != nil {
		query = query.Where("type = ?", *mediaType)
	}

	if err := query.Model(&model.Media{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Preload("Seasons").Find(&mediaList).Error
	if err != nil {
		return nil, 0, err
	}

	return mediaList, total, nil
}

func (r *MediaRepository) Update(ctx context.Context, media *model.Media) error {
	return r.db.WithContext(ctx).Save(media).Error
}

func (r *MediaRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Media{}, id).Error
}

func (r *MediaRepository) Search(ctx context.Context, keyword string) ([]model.Media, error) {
	var mediaList []model.Media
	err := r.db.WithContext(ctx).Where("title LIKE ?", "%"+keyword+"%").Limit(20).Find(&mediaList).Error
	return mediaList, err
}
