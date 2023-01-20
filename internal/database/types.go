package database

import (
	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
)

// GENRE
type GenreSchema struct {
	GenreID   uint32 `gorm:"column:genre_id;primaryKey"`
	Name      string `gorm:"column:name"`
	CreatedAt uint64 `gorm:"column:created_at"`
	UpdatedAt uint64 `gorm:"column:updated_at"`
}

// SONG
// 12 fields
// Put struct
type SongRow struct {
	SongID       uint64
	Name         string
	Genre        common.NameValueInt32
	Artist       string
	Duration     float32
	Language     constants.LANGUAGE_ENUM
	Rating       float32
	ResourceID   string
	ResourceLink string
	CreatedAt    uint64
	UpdatedAt    uint64
	Status       constants.ACTIVE_STATUS `validate:"required"`
}

type SongRowCreate struct {
	Name         string                  `validate:"required,max=256"`
	Genre        common.NameValueInt32   `validate:"required"`
	Artist       string                  `validate:"required"`
	Duration     float32                 `validate:"required"`
	Language     constants.LANGUAGE_ENUM `validate:"required"`
	Rating       float32                 `validate:"required,max=10"`
	ResourceID   string                  `validate:"required"`
	ResourceLink string                  `validate:"required"`
	CreatedAt    uint64
	UpdatedAt    uint64
	Status       constants.ACTIVE_STATUS `validate:"required"`
}

type SongRowPut struct {
	SongID    uint64                  `validate:"required"`
	Name      string                  `validate:"required,max=256"`
	Genre     common.NameValueInt32   `validate:"required"`
	Artist    string                  `validate:"required"`
	Duration  float32                 `validate:"required"`
	Language  constants.LANGUAGE_ENUM `validate:"required"`
	UpdatedAt uint64
	Status    constants.ACTIVE_STATUS `validate:"required"`
}

type UpdateSongResource struct {
	ResourceID   string `validate:"required"`
	ResourceLink string `validate:"required"`
}

type SongSchema struct {
	SongID       uint64                  `gorm:"column:song_id;primaryKey"`
	Name         string                  `gorm:"column:name"`
	Genre        GenreSchema             `gorm:"foreignKey:genre_id"`
	Artist       string                  `gorm:"column:artist"`
	Duration     float32                 `gorm:"column:duration"`
	Language     constants.LANGUAGE_ENUM `gorm:"column:language"`
	Rating       float32                 `gorm:"column:rating"`
	ResourceID   string                  `gorm:"column:resource_id"`
	ResourceLink string                  `gorm:"column:resource_link"`
	CreatedAt    uint64                  `gorm:"column:created_at"`
	UpdatedAt    uint64                  `gorm:"column:updated_at"`
	Status       constants.ACTIVE_STATUS `gorm:"column:status"`
}

// PLAYLIST
type PlaylistSchema struct {
	PlaylistID uint64                  `gorm:"column:playlist_id;primaryKey"`
	Name       string                  `gorm:"column:name"`
	CreatedBy  string                  `gorm:"column:created_by"`
	CreatedAt  uint64                  `gorm:"column:created_at"`
	UpdatedAt  uint64                  `gorm:"column:updated_at"`
	Status     constants.ACTIVE_STATUS `gorm:"column:status"`
	Songs      []SongSchema            `gorm:"many2many:playlist_song;foreignKey:playlist_id;joinForeignKey:playlist_id;References:song_id;joinReferences:song_id"`
}

type PlaylistSchemaCreate struct {
	PlaylistID uint64                  `gorm:"column:playlist_id;primaryKey"`
	Name       string                  `gorm:"column:name" validate:"required,max=256"`
	CreatedBy  string                  `gorm:"column:created_by" validate:"required"`
	CreatedAt  uint64                  `gorm:"column:created_at"`
	UpdatedAt  uint64                  `gorm:"column:updated_at"`
	Status     constants.ACTIVE_STATUS `gorm:"column:status" validate:"required"`
}

type PlaylistSchemaPut struct {
	PlaylistID uint64                  `gorm:"column:playlist_id;primaryKey"`
	Name       string                  `gorm:"column:name" validate:"required,max=256"`
	UpdatedAt  uint64                  `gorm:"column:updated_at"`
	Status     constants.ACTIVE_STATUS `gorm:"column:status" validate:"required"`
}

type DeletePlaylistSchema struct {
	PlaylistID uint64 `gorm:"column:playlist_id;primaryKey"`
}

// PLAYLIST_SONG
type PlaylistSongSchema struct {
	PlaylistID uint64 `gorm:"column:playlist_id;primaryKey"`
	SongID     uint64 `gorm:"column:song_id;primaryKey"`
	CreatedAt  uint64 `gorm:"column:created_at"`
	UpdatedAt  uint64 `gorm:"column:updated_at"`
}

type PlaylistSongSchemaCreate struct {
	PlaylistID uint64 `gorm:"column:playlist_id;primaryKey" validate:"required"`
	SongID     uint64 `gorm:"column:song_id;primaryKey" validate:"required"`
	CreatedAt  uint64 `gorm:"column:created_at"`
	UpdatedAt  uint64 `gorm:"column:updated_at"`
}
