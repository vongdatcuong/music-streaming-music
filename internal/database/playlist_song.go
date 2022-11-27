package database

import (
	"fmt"

	time_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/time"
)

func (db *Database) BatchCreatePlaylistSong(playlistID uint64, songIDs []uint64) ([]PlaylistSongSchemaCreate, error) {
	if songIDs == nil {
		return []PlaylistSongSchemaCreate{}, nil
	}

	playlistSongs := []PlaylistSongSchemaCreate{}

	for _, songID := range songIDs {
		if songID != 0 {
			playlistSongs = append(playlistSongs, PlaylistSongSchemaCreate{PlaylistID: playlistID, SongID: uint64(songID),
				CreatedAt: time_utils.GetCurrentUnixTime(),
				UpdatedAt: time_utils.GetCurrentUnixTime(),
			})
		}
	}

	result := db.GormClient.Create(&playlistSongs)

	if result.Error != nil {
		return []PlaylistSongSchemaCreate{}, fmt.Errorf("could not batch create playlist song records: %w", result.Error)
	}

	return playlistSongs, nil
}
