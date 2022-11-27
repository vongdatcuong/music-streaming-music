package grpc

import (
	"github.com/vongdatcuong/music-streaming-music/internal/modules/playlist"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-music/protos/v1/pb"
)

func convertPlaylistToGrpcPlaylist(myPlaylist playlist.Playlist) *grpcPbV1.Playlist {
	var songs [](*grpcPbV1.Song)

	for _, item := range myPlaylist.Songs {
		songs = append(songs, convertSongToGrpcSong(item))
	}

	return &grpcPbV1.Playlist{
		PlaylistId: myPlaylist.PlaylistID,
		Name:       myPlaylist.Name,
		CreatedBy:  myPlaylist.CreatedBy,
		CreatedAt:  myPlaylist.CreatedAt,
		UpdatedAt:  myPlaylist.UpdatedAt,
		Status:     uint32(myPlaylist.Status),
		Songs:      songs,
	}
}
