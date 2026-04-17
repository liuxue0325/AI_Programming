package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Media    *MediaRepository
	Episode  *EpisodeRepository
	History  *HistoryRepository
	Favorite *FavoriteRepository
	Setting  *SettingRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Media:    NewMediaRepository(db),
		Episode:  NewEpisodeRepository(db),
		History:  NewHistoryRepository(db),
		Favorite: NewFavoriteRepository(db),
		Setting:  NewSettingRepository(db),
	}
}
