package repository

import (
	"context"
	"videohub/backend/internal/model"

	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) Get(ctx context.Context, key string) (*model.Setting, error) {
	var setting model.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *SettingRepository) GetAll(ctx context.Context) ([]model.Setting, error) {
	var settings []model.Setting
	err := r.db.WithContext(ctx).Find(&settings).Error
	return settings, err
}

func (r *SettingRepository) Set(ctx context.Context, key, value, description string) error {
	var setting model.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&setting).Error

	if err == gorm.ErrRecordNotFound {
		setting = model.Setting{
			Key:         key,
			Value:       value,
			Description: description,
		}
		return r.db.WithContext(ctx).Create(&setting).Error
	} else if err != nil {
		return err
	}

	setting.Value = value
	setting.Description = description
	return r.db.WithContext(ctx).Save(&setting).Error
}

func (r *SettingRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key = ?", key).Delete(&model.Setting{}).Error
}
