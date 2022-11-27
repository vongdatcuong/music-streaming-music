package database

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/playlist"
	common_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/common"
	time_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/time"
	validator_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/validator"
)

func (db *Database) GetPlaylistList(ctx context.Context, paginationInfo common.PaginationInfo, filter playlist.PlaylistListFilter) ([]playlist.Playlist, uint64, error) {
	var playlistSchemas []PlaylistSchema
	var totalCount int64

	result := db.GormClient.WithContext(ctx).Preload("Songs.Genre").Preload("Songs").Order("created_at desc, name").Where(PlaylistSchema{Name: filter.Name, CreatedBy: filter.CreatedBy})

	if filter.CreatedTimeFrom != 0 {
		result = result.Where("created_at >= ?", filter.CreatedTimeFrom)
	}

	if filter.CreatedTimeTo != 0 {
		result = result.Where("created_at <= ?", filter.CreatedTimeTo)
	}

	result = result.Scopes(Paginate(paginationInfo)).Find(&playlistSchemas).Count(&totalCount)

	if result.Error != nil {
		return []playlist.Playlist{}, 0, fmt.Errorf("could not get the playlist list: %w", result.Error)
	}

	var playlists []playlist.Playlist

	for _, model := range playlistSchemas {
		playlists = append(playlists, convertPlaylistSchemaToPlaylist(model))
	}

	return playlists, uint64(totalCount), nil
}

func (db *Database) GetPlaylistDetails(ctx context.Context, id uint64) (playlist.Playlist, error) {
	var record PlaylistSchema

	result := db.GormClient.WithContext(ctx).Preload("Songs.Genre").Preload("Songs").First(&record, id)
	if result.Error != nil {
		return playlist.Playlist{}, fmt.Errorf("could not get playlist details: %w", result.Error)
	}

	return convertPlaylistSchemaToPlaylist(record), nil
}

func (db *Database) CreatePlaylist(ctx context.Context, newPlaylist playlist.Playlist, songIDs []uint64) (playlist.Playlist, error) {
	playlistRowCreate := PlaylistSchemaCreate{
		Name:      newPlaylist.Name,
		CreatedBy: newPlaylist.CreatedBy,
		CreatedAt: time_utils.GetCurrentUnixTime(),
		UpdatedAt: time_utils.GetCurrentUnixTime(),
		Status:    newPlaylist.Status,
	}

	err := validator_utils.ValidateStruct(playlistRowCreate)

	if err != nil {
		return playlist.Playlist{}, fmt.Errorf("playlist is not valid: %w", err)
	}

	result := db.GormClient.WithContext(ctx).Create(&playlistRowCreate)

	if result.Error != nil {
		return playlist.Playlist{}, fmt.Errorf("could not create playlist %w: ", result.Error)
	}

	_, err = db.BatchCreatePlaylistSong(playlistRowCreate.PlaylistID, songIDs)

	if err != nil {
		return playlist.Playlist{}, err
	}

	return convertPlaylistSchemaCreateToPlaylist(playlistRowCreate), nil
}

func (db *Database) PutPlaylist(ctx context.Context, curPlaylist playlist.Playlist) (playlist.Playlist, error) {
	playlistRowPut := PlaylistSchemaPut{
		PlaylistID: curPlaylist.PlaylistID,
		Name:       curPlaylist.Name,
		UpdatedAt:  time_utils.GetCurrentUnixTime(),
		Status:     curPlaylist.Status,
	}

	err := validator_utils.ValidateStruct(playlistRowPut)

	if err != nil {
		return playlist.Playlist{}, fmt.Errorf("playlist is not valid: %w", err)
	}

	result := db.GormClient.WithContext(ctx).Updates(playlistRowPut)

	if result.Error != nil {
		return playlist.Playlist{}, fmt.Errorf("could not put playlist %w: ", result.Error)
	}

	return curPlaylist, nil
}

func (db *Database) DeletePlaylist(ctx context.Context, id uint64) error {
	result := db.GormClient.WithContext(ctx).Delete(&DeletePlaylistSchema{PlaylistID: id})

	if result.Error != nil {
		return fmt.Errorf("could not delete playlist: %w", result.Error)
	}

	return nil
}

func (db *Database) DoesPlaylistExist(ctx context.Context, id uint64) (bool, error) {
	var exists bool

	result := db.GormClient.WithContext(ctx).Table(PlaylistTableName).Select("count(*) > 0").Where("playlist_id = ?", id).Find(&exists)

	if result.Error != nil {
		return false, fmt.Errorf("could not check if playlist exists: %w", result.Error)
	}

	return exists, nil
}

func (db *Database) UpdatePlaylistSongs(ctx context.Context, playlistID uint64, songIDs []uint64) error {
	curSongIDs, err := db.GetSongsOfAPlaylist(ctx, playlistID)

	if err != nil {
		return err
	}

	var deletedSongIDs []uint64
	var newSongIDs []uint64

	for _, id := range songIDs {
		if !common_utils.Contains(curSongIDs, id) {
			newSongIDs = append(newSongIDs, id)
		}
	}

	for _, id := range curSongIDs {
		if !common_utils.Contains(songIDs, id) {
			deletedSongIDs = append(deletedSongIDs, id)
		}
	}
	logrus.Info(newSongIDs)
	logrus.Info(deletedSongIDs)
	err = db.AddSongsToAPlaylist(ctx, playlistID, newSongIDs)

	if err != nil {
		return err
	}

	err = db.DeleteSongsFromAPlaylist(ctx, playlistID, deletedSongIDs)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetSongsOfAPlaylist(ctx context.Context, playlistID uint64) ([]uint64, error) {
	var songIDs []uint64

	result := db.GormClient.WithContext(ctx).Model(&PlaylistSongSchema{}).Where("playlist_id = ?", playlistID).Pluck("song_id", &songIDs)

	if result.Error != nil {
		return []uint64{}, fmt.Errorf("could not get all songs of the playlist %d: %w", playlistID, result.Error)
	}

	return songIDs, nil
}

func (db *Database) AddSongsToAPlaylist(ctx context.Context, playlistID uint64, songIDs []uint64) error {
	if len(songIDs) == 0 {
		return nil
	}

	var playlistSongs []PlaylistSongSchema

	for _, songID := range songIDs {
		playlistSongs = append(playlistSongs, PlaylistSongSchema{PlaylistID: playlistID, SongID: songID,
			CreatedAt: time_utils.GetCurrentUnixTime(), UpdatedAt: time_utils.GetCurrentUnixTime()})
	}

	result := db.GormClient.WithContext(ctx).Model(&PlaylistSongSchema{}).Create(&playlistSongs)

	if result.Error != nil {
		return fmt.Errorf("could not add songs to the playlist %d: %w", playlistID, result.Error)
	}

	return nil
}

func (db *Database) DeleteSongsFromAPlaylist(ctx context.Context, playlistID uint64, songIDs []uint64) error {
	if len(songIDs) == 0 {
		return nil
	}

	var playlistSongs []PlaylistSongSchema

	for _, songID := range songIDs {
		playlistSongs = append(playlistSongs, PlaylistSongSchema{PlaylistID: playlistID, SongID: songID})
	}

	result := db.GormClient.WithContext(ctx).Delete(playlistSongs)

	if result.Error != nil {
		return fmt.Errorf("could not delete songs of the playlist %d: %w", playlistID, result.Error)
	}

	return nil
}
