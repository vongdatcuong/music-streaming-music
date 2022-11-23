package grpc

import (
	"context"
	"mime/multipart"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
	common_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-music/protos/v1/pb"
)

type SongServiceGrpc interface {
	GetSongList(context.Context, common.PaginationInfo, song.SongListFilter) ([]song.Song, uint64, error)
	GetSongDetails(context.Context, uint64) (song.Song, error)
	CreateSong(context.Context, song.Song) (song.Song, error)
	PutSong(context.Context, song.Song) (song.Song, error)
	DeleteSong(context.Context, uint64) error
	UploadSong(context.Context, *multipart.FileHeader, multipart.File) (string, string, error)
}

func convertSongToGrpcSong(mySong song.Song) *grpcPbV1.Song {
	return &grpcPbV1.Song{
		SongId:       mySong.SongID,
		Name:         mySong.Name,
		Genre:        &grpcPbV1.NameValueInt32{Name: mySong.Genre.Name, Value: mySong.Genre.Value},
		Artist:       mySong.Artist,
		Duration:     mySong.Duration,
		Language:     string(mySong.Language),
		Rating:       mySong.Rating,
		ResourceId:   mySong.ResourceID,
		ResourceLink: mySong.ResourceLink,
		CreatedAt:    mySong.CreatedAt,
		UpdatedAt:    mySong.UpdatedAt,
		Status:       uint32(mySong.Status),
	}
}

func (h *Handler) GetSongList(ctx context.Context, req *grpcPbV1.GetSongListRequest) (*grpcPbV1.GetSongListResponse, error) {
	var songPagination common.PaginationInfo = common.PaginationInfo{}
	var songListFilter song.SongListFilter = song.SongListFilter{}

	if req.PaginationInfo != nil {
		songPagination = common.PaginationInfo{
			Offset: req.PaginationInfo.Offset,
			Limit:  req.PaginationInfo.Limit,
		}
	}

	if req.Filter != nil {
		if req.Filter.Name != "" {
			songListFilter.Name = &req.Filter.Name
		}
		if req.Filter.Genre != 0 {
			songListFilter.Genre = &req.Filter.Genre
		}
		if req.Filter.Artist != "" {
			songListFilter.Artist = &req.Filter.Artist
		}
		if req.Filter.Duration != 0 {
			songListFilter.Duration = &req.Filter.Duration
		}
		if req.Filter.Language != "" {
			language := constants.LANGUAGE_ENUM(req.Filter.Language)
			songListFilter.Language = &language
		}
		if req.Filter.CreatedTimeFrom != 0 {
			songListFilter.CreatedTimeFrom = &req.Filter.CreatedTimeFrom
		}
		if req.Filter.CreatedTimeTo != 0 {
			songListFilter.CreatedTimeTo = &req.Filter.CreatedTimeTo
		}
	}

	songs, totalCount, err := h.songService.GetSongList(ctx, songPagination, songListFilter)

	if err != nil {
		return &grpcPbV1.GetSongListResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	grpcSongs := [](*grpcPbV1.Song){}

	for i := 0; i < len(songs); i++ {
		grpcSongs = append(grpcSongs, convertSongToGrpcSong(songs[i]))
	}
	return &grpcPbV1.GetSongListResponse{
		Data: &grpcPbV1.GetSongListResponseData{
			Songs:      grpcSongs,
			TotalCount: &totalCount,
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) GetSongDetails(ctx context.Context, req *grpcPbV1.GetSongDetailsRequest) (*grpcPbV1.GetSongDetailsResponse, error) {
	fetchedSong, err := h.songService.GetSongDetails(ctx, req.SongId)

	if err != nil {
		return &grpcPbV1.GetSongDetailsResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.GetSongDetailsResponse{
		Data: &grpcPbV1.GetSongDetailsResponseData{
			Song: convertSongToGrpcSong(fetchedSong),
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) CreateSong(ctx context.Context, req *grpcPbV1.CreateSongRequest) (*grpcPbV1.CreateSongResponse, error) {
	newSong := song.Song{
		Name:         req.Song.Name,
		Genre:        common.NameValueInt32{Name: "", Value: req.Song.Genre.Value},
		Artist:       req.Song.Artist,
		Language:     constants.LANGUAGE_ENUM(req.Song.Language),
		Duration:     req.Song.Duration, // Passed by FE
		ResourceID:   req.Song.ResourceId,
		ResourceLink: req.Song.ResourceLink,
	}
	_, err := h.songService.CreateSong(ctx, newSong)

	if err != nil {
		return &grpcPbV1.CreateSongResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.CreateSongResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) PutSong(ctx context.Context, req *grpcPbV1.PutSongRequest) (*grpcPbV1.PutSongResponse, error) {
	newSong := song.Song{
		SongID:   req.Song.SongId,
		Name:     req.Song.Name,
		Genre:    common.NameValueInt32{Name: "", Value: req.Song.Genre.Value},
		Artist:   req.Song.Artist,
		Duration: req.Song.Duration,
		Language: constants.LANGUAGE_ENUM(req.Song.Language),
		Status:   constants.ACTIVE_STATUS(req.Song.Status),
	}
	_, err := h.songService.PutSong(ctx, newSong)

	if err != nil {
		return &grpcPbV1.PutSongResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.PutSongResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) DeleteSong(ctx context.Context, req *grpcPbV1.DeleteSongRequest) (*grpcPbV1.DeleteSongResponse, error) {
	err := h.songService.DeleteSong(ctx, req.SongId)

	if err != nil {
		return &grpcPbV1.DeleteSongResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.DeleteSongResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}
