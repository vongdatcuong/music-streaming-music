package grpc

import (
	"context"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/playlist"
	common_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-music/protos/v1/pb"
)

type PlaylistServiceGrpc interface {
	GetPlaylistList(context.Context, common.PaginationInfo, playlist.PlaylistListFilter) ([]playlist.Playlist, uint64, error)
	GetPlaylistDetails(context.Context, uint64) (playlist.Playlist, error)
	CreatePlaylist(context.Context, playlist.Playlist) (playlist.Playlist, error)
	PutPlaylist(context.Context, playlist.Playlist) (playlist.Playlist, error)
	DeletePlaylist(context.Context, uint64) error
}

func convertPlaylistToGrpcPlaylist(myPlaylist playlist.Playlist) *grpcPbV1.Playlist {
	return &grpcPbV1.Playlist{
		PlaylistId: myPlaylist.PlaylistID,
		Name:       myPlaylist.Name,
		CreatedBy: &grpcPbV1.UserDetail{
			UserId: myPlaylist.CreatedBy.UserID,
			Email:  myPlaylist.CreatedBy.Email,
		},
		CreatedAt: myPlaylist.CreatedAt,
		UpdatedAt: myPlaylist.UpdatedAt,
		Status:    uint32(myPlaylist.Status),
	}
}

func (h *Handler) GetPlaylistList(ctx context.Context, req *grpcPbV1.GetPlaylistListRequest) (*grpcPbV1.GetPlaylistListResponse, error) {
	var pagination common.PaginationInfo = common.PaginationInfo{}
	var playlistListFilter playlist.PlaylistListFilter = playlist.PlaylistListFilter{}

	if req.PaginationInfo != nil {
		pagination = common.PaginationInfo{
			Offset: req.PaginationInfo.Offset,
			Limit:  req.PaginationInfo.Limit,
		}
	}

	if req.Filter != nil {
		if req.Filter.Name != "" {
			playlistListFilter.Name = req.Filter.Name
		}
		if req.Filter.CreatedBy != "" {
			playlistListFilter.CreatedBy = req.Filter.CreatedBy
		}
		if req.Filter.CreatedTimeFrom != 0 {
			playlistListFilter.CreatedTimeFrom = req.Filter.CreatedTimeFrom
		}
		if req.Filter.CreatedTimeTo != 0 {
			playlistListFilter.CreatedTimeTo = req.Filter.CreatedTimeTo
		}
	}

	playlists, totalCount, err := h.playlistService.GetPlaylistList(ctx, pagination, playlistListFilter)

	if err != nil {
		return &grpcPbV1.GetPlaylistListResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	grpcPlaylists := [](*grpcPbV1.Playlist){}

	for i := 0; i < len(playlists); i++ {
		grpcPlaylists = append(grpcPlaylists, convertPlaylistToGrpcPlaylist(playlists[i]))
	}
	return &grpcPbV1.GetPlaylistListResponse{
		Data: &grpcPbV1.GetPlaylistListResponseData{
			Playlists:  grpcPlaylists,
			TotalCount: &totalCount,
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) GetPlaylistDetails(ctx context.Context, req *grpcPbV1.GetPlaylistDetailsRequest) (*grpcPbV1.GetPlaylistDetailsResponse, error) {
	fetchedPlaylist, err := h.playlistService.GetPlaylistDetails(ctx, req.PlaylistId)

	if err != nil {
		return &grpcPbV1.GetPlaylistDetailsResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.GetPlaylistDetailsResponse{
		Data: &grpcPbV1.GetPlaylistDetailsResponseData{
			Playlist: convertPlaylistToGrpcPlaylist(fetchedPlaylist),
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) CreatePlaylist(ctx context.Context, req *grpcPbV1.CreatePlaylistRequest) (*grpcPbV1.CreatePlaylistResponse, error) {
	newPlaylist := playlist.Playlist{
		Name: req.Playlist.Name,
		CreatedBy: common.UserDetail{
			UserID: req.Playlist.CreatedBy.UserId,
		},
	}

	_, err := h.playlistService.CreatePlaylist(ctx, newPlaylist)

	if err != nil {
		return &grpcPbV1.CreatePlaylistResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.CreatePlaylistResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) PutPlaylist(ctx context.Context, req *grpcPbV1.PutPlaylistRequest) (*grpcPbV1.PutPlaylistResponse, error) {
	newPlaylist := playlist.Playlist{
		PlaylistID: req.Playlist.PlaylistId,
		Name:       req.Playlist.Name,
		CreatedBy: common.UserDetail{
			UserID: req.Playlist.CreatedBy.UserId,
		},
		Status: constants.ACTIVE_STATUS(req.Playlist.Status),
	}
	_, err := h.playlistService.PutPlaylist(ctx, newPlaylist)

	if err != nil {
		return &grpcPbV1.PutPlaylistResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.PutPlaylistResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) DeletePlaylist(ctx context.Context, req *grpcPbV1.DeletePlaylistRequest) (*grpcPbV1.DeletePlaylistResponse, error) {
	err := h.playlistService.DeletePlaylist(ctx, req.PlaylistId)

	if err != nil {
		return &grpcPbV1.DeletePlaylistResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.DeletePlaylistResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}