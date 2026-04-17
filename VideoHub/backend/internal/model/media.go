package model

import (
	"time"
)

type MediaType string

const (
	MediaTypeMovie MediaType = "movie"
	MediaTypeTV    MediaType = "tv"
	MediaTypeAnime MediaType = "anime"
)

type MediaStatus string

const (
	MediaStatusPending   MediaStatus = "pending"
	MediaStatusScraped   MediaStatus = "scraped"
	MediaStatusFailed    MediaStatus = "failed"
)

type Media struct {
	ID            int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Type          MediaType   `gorm:"type:text;not null;index" json:"type"`
	Title         string      `gorm:"type:text;not null" json:"title"`
	OriginalTitle string      `gorm:"type:text" json:"original_title,omitempty"`
	Year          int         `gorm:"type:integer;index" json:"year,omitempty"`
	Path          string      `gorm:"type:text;not null;uniqueIndex" json:"path"`
	PosterPath    string      `gorm:"type:text" json:"poster_path,omitempty"`
	BackdropPath  string      `gorm:"type:text" json:"backdrop_path,omitempty"`
	Overview      string      `gorm:"type:text" json:"overview,omitempty"`
	Rating        float64     `gorm:"type:real;default:0;index" json:"rating,omitempty"`
	Runtime       int         `gorm:"type:integer;default:0" json:"runtime,omitempty"`
	Genres        string      `gorm:"type:text" json:"genres,omitempty"`
	TmdbID        int64       `gorm:"type:integer;index" json:"tmdb_id,omitempty"`
	ImdbID        string      `gorm:"type:text" json:"imdb_id,omitempty"`
	Status        MediaStatus `gorm:"type:text;default:'pending';index" json:"status"`
	Seasons       []Season    `gorm:"foreignKey:SeriesID" json:"seasons,omitempty"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type Season struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SeriesID     int64      `gorm:"not null;index" json:"series_id"`
	SeasonNumber int        `gorm:"not null" json:"season_number"`
	Title        string     `gorm:"type:text" json:"title,omitempty"`
	PosterPath   string     `gorm:"type:text" json:"poster_path,omitempty"`
	Episodes     []Episode  `gorm:"foreignKey:SeriesID,SeasonNumber" json:"episodes,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Episode struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SeriesID      int64     `gorm:"not null;index:idx_series_id" json:"series_id"`
	SeasonNumber  int       `gorm:"not null" json:"season_number"`
	EpisodeNumber int       `gorm:"not null" json:"episode_number"`
	Title         string    `gorm:"type:text" json:"title,omitempty"`
	Path          string    `gorm:"type:text;not null;uniqueIndex" json:"path"`
	PosterPath    string    `gorm:"type:text" json:"poster_path,omitempty"`
	Overview      string    `gorm:"type:text" json:"overview,omitempty"`
	Runtime       int       `gorm:"type:integer;default:0" json:"runtime,omitempty"`
	TmdbID        int64     `gorm:"type:integer" json:"tmdb_id,omitempty"`
	AirDate       string    `gorm:"type:text" json:"air_date,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
