package database

import (
	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
)

func convertSongRowToSong(songRow SongRow) song.Song {
	return song.Song{
		SongID:       songRow.SongID,
		Name:         songRow.Name,
		Genre:        songRow.Genre,
		Artist:       songRow.Artist,
		Duration:     songRow.Duration,
		Language:     songRow.Language,
		Rating:       songRow.Rating,
		ResourceID:   songRow.ResourceID,
		ResourceLink: songRow.ResourceLink,
		CreatedAt:    songRow.CreatedAt,
		UpdatedAt:    songRow.UpdatedAt,
		Status:       songRow.Status,
	}
}

func convertSongSchemaToSong(schema SongSchema) song.Song {
	return song.Song{
		SongID: schema.SongID,
		Name:   schema.Name,
		Genre: common.NameValueInt32{
			Name:  schema.Genre.Name,
			Value: int32(schema.Genre.GenreID),
		},
		Artist:       schema.Artist,
		Duration:     schema.Duration,
		Language:     schema.Language,
		Rating:       schema.Rating,
		ResourceID:   schema.ResourceID,
		ResourceLink: schema.ResourceLink,
		CreatedAt:    schema.CreatedAt,
		UpdatedAt:    schema.UpdatedAt,
		Status:       schema.Status,
	}
}
