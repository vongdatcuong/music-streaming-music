package song

import (
	"context"
	"fmt"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	constants "github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
)

type SongListFilter struct {
	Name            *string
	Genre           *uint32
	Artist          *string
	Duration        *uint32
	Language        *constants.LANGUAGE_ENUM
	CreatedTimeFrom *uint64
	CreatedTimeTo   *uint64
}

type SongStore interface {
	GetSongList(context.Context, common.PaginationInfo, SongListFilter) ([]Song, uint64, error)
	GetSongDetails(context.Context, uint64) (Song, error)
	CreateSong(context.Context, Song) (Song, error)
	PutSong(context.Context, Song) (Song, error)
	DeleteSong(context.Context, uint64) error
	DoesSongExist(context.Context, uint64) (bool, error)
}

type SongService struct {
	store SongStore
}

type Song struct {
	SongID       uint64
	Name         string
	Genre        common.NameValueInt32
	Artist       string
	Duration     uint32
	Language     constants.LANGUAGE_ENUM
	Rating       float32
	ResourceID   uint64
	ResourceLink string
	CreatedAt    uint64
	UpdatedAt    uint64
	Status       constants.ACTIVE_STATUS
}

func NewService(store SongStore) *SongService {
	return &SongService{
		store: store,
	}
}

func (s *SongService) GetSongList(ctx context.Context, pagination common.PaginationInfo, filter SongListFilter) ([]Song, uint64, error) {
	songList, totalCount, err := s.store.GetSongList(ctx, pagination, filter)

	if err != nil {
		return []Song{}, 0, err
	}

	return songList, totalCount, nil
}

func (s *SongService) GetSongDetails(ctx context.Context, id uint64) (Song, error) {
	song, err := s.store.GetSongDetails(ctx, id)

	if err != nil {
		return Song{}, err
	}

	return song, nil
}

func (s *SongService) CreateSong(ctx context.Context, newSong Song) (Song, error) {
	wrappedSong := Song(newSong)
	wrappedSong.Duration = 234 // calculate duration from the audio file its self
	wrappedSong.ResourceID = 1
	wrappedSong.ResourceLink = "empty link"
	wrappedSong.Rating = 11
	wrappedSong.Status = constants.ACTIVE_STATUS_ACTIVE

	song, err := s.store.CreateSong(ctx, wrappedSong)

	if err != nil {
		return Song{}, err
	}

	return song, nil
}

func (s *SongService) PutSong(ctx context.Context, existingSong Song) (Song, error) {
	songID := existingSong.SongID
	doesExist, err := s.store.DoesSongExist(ctx, songID)

	if err != nil {
		return Song{}, nil
	}

	if !doesExist {
		return Song{}, fmt.Errorf("could not find song with id %d", songID)
	}

	song, err := s.store.PutSong(ctx, existingSong)

	if err != nil {
		return Song{}, err
	}

	return song, nil
}

func (s *SongService) DeleteSong(ctx context.Context, id uint64) error {
	doesExist, err := s.store.DoesSongExist(ctx, id)

	if err != nil {
		return nil
	}

	if !doesExist {
		return fmt.Errorf("could not find song with id %d", id)
	}

	err = s.store.DeleteSong(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
