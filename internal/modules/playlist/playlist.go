package playlist

import (
	"context"
	"fmt"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	constants "github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
)

type PlaylistListFilter struct {
	Name            string
	CreatedBy       string
	CreatedTimeFrom uint64
	CreatedTimeTo   uint64
}

type PlaylistStore interface {
	GetPlaylistList(context.Context, common.PaginationInfo, PlaylistListFilter) ([]Playlist, uint64, error)
	GetPlaylistDetails(context.Context, uint64) (Playlist, error)
	CreatePlaylist(context.Context, Playlist, []uint64) (Playlist, error)
	PutPlaylist(context.Context, Playlist) (Playlist, error)
	DeletePlaylist(context.Context, uint64) error
	DoesPlaylistExist(context.Context, uint64) (bool, error)
	UpdatePlaylistSongs(context.Context, uint64, []uint64) error
}

type PlaylistService struct {
	store PlaylistStore
}

type Playlist struct {
	PlaylistID uint64
	Name       string
	CreatedBy  string
	CreatedAt  uint64
	UpdatedAt  uint64
	Status     constants.ACTIVE_STATUS
	Songs      []song.Song
}

func NewService(store PlaylistStore) *PlaylistService {
	return &PlaylistService{
		store: store,
	}
}

func (s *PlaylistService) GetPlaylistList(ctx context.Context, pagination common.PaginationInfo, filter PlaylistListFilter) ([]Playlist, uint64, error) {
	playlistList, totalCount, err := s.store.GetPlaylistList(ctx, pagination, filter)

	if err != nil {
		return []Playlist{}, 0, err
	}

	return playlistList, totalCount, nil
}

func (s *PlaylistService) GetPlaylistDetails(ctx context.Context, id uint64) (Playlist, error) {
	playlist, err := s.store.GetPlaylistDetails(ctx, id)

	if err != nil {
		return Playlist{}, err
	}

	return playlist, nil
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, newPlaylist Playlist, songIDs []uint64) (Playlist, error) {
	wrappedPlaylist := Playlist(newPlaylist)
	wrappedPlaylist.Status = constants.ACTIVE_STATUS_ACTIVE
	playlist, err := s.store.CreatePlaylist(ctx, wrappedPlaylist, songIDs)

	if err != nil {
		return Playlist{}, err
	}

	return playlist, nil
}

func (s *PlaylistService) PutPlaylist(ctx context.Context, existingPlaylist Playlist) (Playlist, error) {
	playlistID := existingPlaylist.PlaylistID
	doesExist, err := s.store.DoesPlaylistExist(ctx, playlistID)

	if err != nil {
		return Playlist{}, err
	}

	if !doesExist {
		return Playlist{}, fmt.Errorf("could not find playlist with id %d", playlistID)
	}

	playlist, err := s.store.PutPlaylist(ctx, existingPlaylist)

	if err != nil {
		return Playlist{}, err
	}

	return playlist, nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id uint64) error {
	doesExist, err := s.store.DoesPlaylistExist(ctx, id)

	if err != nil {
		return err
	}

	if !doesExist {
		return fmt.Errorf("could not find playlist with id %d", id)
	}

	err = s.store.DeletePlaylist(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistService) UpdatePlaylistSongs(ctx context.Context, playlistID uint64, songIDs []uint64) error {
	doesExist, err := s.store.DoesPlaylistExist(ctx, playlistID)

	if err != nil {
		return err
	}

	if !doesExist {
		return fmt.Errorf("could not find playlist with id %d", playlistID)
	}

	err = s.store.UpdatePlaylistSongs(ctx, playlistID, songIDs)

	if err != nil {
		return err
	}

	return nil
}
