package model

import (
	"time"
)

type Favorite struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MediaID   int64     `gorm:"not null;uniqueIndex" json:"media_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Media     *Media    `gorm:"foreignKey:MediaID" json:"media,omitempty"`
}
