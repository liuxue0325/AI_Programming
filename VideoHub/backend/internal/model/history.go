package model

import (
	"time"
)

type WatchHistory struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	MediaID     int64      `gorm:"index" json:"media_id"`
	EpisodeID   *int64     `gorm:"index" json:"episode_id,omitempty"`
	Progress    float64    `gorm:"type:real;default:0" json:"progress"`
	Completed   bool       `gorm:"type:boolean;default:false" json:"completed"`
	LastWatched time.Time  `gorm:"autoUpdateTime;index" json:"last_watched"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Media       *Media     `gorm:"foreignKey:MediaID" json:"media,omitempty"`
	Episode     *Episode   `gorm:"foreignKey:EpisodeID" json:"episode,omitempty"`
}
