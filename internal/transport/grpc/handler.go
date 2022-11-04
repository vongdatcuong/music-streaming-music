package grpc

import (
	"fmt"
	"net"

	grpcPbV1 "github.com/vongdatcuong/music-streaming-music/protos/v1/pb"
	"google.golang.org/grpc"
)

type Handler struct {
	grpcPbV1.UnimplementedSongServiceServer
	grpcPbV1.UnimplementedPlaylistServiceServer
}

func NewHandler() *Handler {
	h := &Handler{}

	return h
}

func (h *Handler) Server() error {
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		return fmt.Errorf("could not listen on port 8001: %w", err)
	}

	grpcServer := grpc.NewServer()
	grpcPbV1.RegisterSongServiceServer(grpcServer, h)
	grpcPbV1.RegisterPlaylistServiceServer(grpcServer, h)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("could not server Grpc server on port 8001: %w", err)
	}

	return nil
}
