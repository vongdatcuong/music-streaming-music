package genre

import (
	"context"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
)

type GenreStore interface {
	GetGenreOptionsList(context.Context) ([]common.NameValueInt32, error)
}

type GenreService struct {
	store GenreStore
}

func NewService(store GenreStore) *GenreService {
	return &GenreService{
		store: store,
	}
}

func (s *GenreService) GetGenreOptionsList(ctx context.Context) ([]common.NameValueInt32, error) {
	optionsList, err := s.store.GetGenreOptionsList(ctx)

	if err != nil {
		return []common.NameValueInt32{}, err
	}

	return optionsList, nil
}
