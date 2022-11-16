package database

import (
	"context"
	"fmt"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/playlist"
	time_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/time"
	validator_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/validator"
)

type PlaylistRowCreate struct {
	PlaylistID uint64
	Name       string `validate:"required,max=256"`
	CreatedBy  uint64 `validate:"required"`
	CreatedAt  uint64
	UpdatedAt  uint64
	Status     constants.ACTIVE_STATUS `validate:"required"`
}

func convertPlaylistRowCreateToPlaylist(playlistRowCreate PlaylistRowCreate) playlist.Playlist {
	return playlist.Playlist{
		PlaylistID: playlistRowCreate.PlaylistID,
		Name:       playlistRowCreate.Name,
		CreatedBy: common.UserDetail{
			UserID: playlistRowCreate.CreatedBy,
			Email:  "temp@gmail.com", // TODO: Update this when building User service
		},
		CreatedAt: playlistRowCreate.CreatedAt,
		UpdatedAt: playlistRowCreate.UpdatedAt,
		Status:    playlistRowCreate.Status,
	}
}

func (db *Database) CreatePlaylist(ctx context.Context, newPlaylist playlist.Playlist) (playlist.Playlist, error) {
	playlistRowCreate := PlaylistRowCreate{
		Name:      newPlaylist.Name,
		CreatedBy: newPlaylist.CreatedBy.UserID,
		CreatedAt: time_utils.GetCurrentUnixTime(),
		UpdatedAt: time_utils.GetCurrentUnixTime(),
		Status:    newPlaylist.Status,
	}

	err := validator_utils.ValidateStruct(playlistRowCreate)

	if err != nil {
		return playlist.Playlist{}, fmt.Errorf("playlist is not valid: %w", err)
	}

	result := db.GormClient.Table("Playlist").Create(&playlistRowCreate)

	if result.Error != nil {
		return playlist.Playlist{}, fmt.Errorf("could not create playlist %w: ", result.Error)
	}
	return convertPlaylistRowCreateToPlaylist(playlistRowCreate), nil
}
