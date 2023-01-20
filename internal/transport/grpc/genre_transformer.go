package grpc

import (
	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
)

func convertNameInt32ToGrpcNameInt32Value(pair common.NameValueInt32) *grpcPbV1.NameValueInt32 {
	return &grpcPbV1.NameValueInt32{
		Name:  pair.Name,
		Value: pair.Value,
	}
}
