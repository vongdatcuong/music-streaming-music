package database

import "github.com/vongdatcuong/music-streaming-music/internal/modules/song"

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
