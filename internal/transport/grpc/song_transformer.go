package grpc

import (
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
)

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
