package grpc

import (
	"context"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	common_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
)

type GenreServiceGrpc interface {
	GetGenreOptionsList(context.Context) ([]common.NameValueInt32, error)
}

func (h *Handler) GetGenreOptionsList(ctx context.Context, req *grpcPbV1.GetGenreOptionsListRequest) (*grpcPbV1.GetGenreOptionsListResponse, error) {
	optionsList, err := h.genreService.GetGenreOptionsList(ctx)

	if err != nil {
		return &grpcPbV1.GetGenreOptionsListResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	var convertedOptionsList [](*grpcPbV1.NameValueInt32)

	for _, option := range optionsList {
		convertedOptionsList = append(convertedOptionsList, convertNameInt32ToGrpcNameInt32Value(option))
	}

	return &grpcPbV1.GetGenreOptionsListResponse{
		Data: &grpcPbV1.GetGenreOptionsListResponseData{
			Genres: convertedOptionsList,
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}
