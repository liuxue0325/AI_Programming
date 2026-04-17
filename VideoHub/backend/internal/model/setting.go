package model

import (
	"time"
)

type Setting struct {
	Key         string    `gorm:"primaryKey;type:text" json:"key"`
	Value       string    `gorm:"type:text" json:"value"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ScanTask struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FolderPath     string    `gorm:"type:text;not null;index" json:"folder_path"`
	Status         string    `gorm:"type:text;default:'pending';index" json:"status"`
	TotalFiles     int       `gorm:"type:integer;default:0" json:"total_files"`
	ProcessedFiles int       `gorm:"type:integer;default:0" json:"processed_files"`
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
	CompletedAt    time.Time `json:"completed_at,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}
