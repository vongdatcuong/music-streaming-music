package database

import (
	"github.com/vongdatcuong/music-streaming-music/internal/modules/playlist"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
)

func convertPlaylistSchemaToPlaylist(schema PlaylistSchema) playlist.Playlist {
	var songs []song.Song

	for _, item := range schema.Songs {
		songs = append(songs, convertSongSchemaToSong(item))
	}

	return playlist.Playlist{
		PlaylistID: schema.PlaylistID,
		Name:       schema.Name,
		CreatedBy:  schema.CreatedBy,
		CreatedAt:  schema.CreatedAt,
		UpdatedAt:  schema.UpdatedAt,
		Status:     schema.Status,
		Songs:      songs,
	}
}

func convertPlaylistSchemaCreateToPlaylist(schema PlaylistSchemaCreate) playlist.Playlist {
	return playlist.Playlist{
		PlaylistID: schema.PlaylistID,
		Name:       schema.Name,
		CreatedBy:  schema.CreatedBy,
		CreatedAt:  schema.CreatedAt,
		UpdatedAt:  schema.UpdatedAt,
		Status:     schema.Status,
	}
}
