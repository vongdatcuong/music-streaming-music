package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-music/protos/v1/pb"
	"google.golang.org/grpc"
)

type Handler struct {
	grpcPbV1.UnimplementedSongServiceServer
	grpcPbV1.UnimplementedPlaylistServiceServer
	songService     SongServiceGrpc
	playlistService PlaylistServiceGrpc
	storageService  StorageService
}

type StorageService interface {
	CreateFile(content []byte, fileName string, folderPath string) (string, string, error)
}

func NewHandler(songService SongServiceGrpc, storageService StorageService) *Handler {
	h := &Handler{songService: songService, storageService: storageService}

	return h
}

func (h *Handler) RunGrpcServer(port string, channel chan error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		channel <- fmt.Errorf("could not listen on port %s: %w", port, err)
		return
	}

	grpcServer := grpc.NewServer()
	grpcPbV1.RegisterSongServiceServer(grpcServer, h)
	grpcPbV1.RegisterPlaylistServiceServer(grpcServer, h)

	if err := grpcServer.Serve(lis); err != nil {
		channel <- fmt.Errorf("could not server Grpc server on port %s: %w", port, err)
	}
}

func (h *Handler) RunRestServer(port string, channel chan error) {
	gwmux := runtime.NewServeMux()
	muxCtx, cancelMuxCtx := context.WithCancel(context.Background())
	defer cancelMuxCtx()
	err := grpcPbV1.RegisterSongServiceHandlerServer(muxCtx, gwmux, h)

	if err != nil {
		channel <- fmt.Errorf("Failed to register gateway: %w", err)
		return
	}

	restLis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		channel <- fmt.Errorf("could not listen on port %s: %w", port, err)
		return
	}
	httpMux := http.NewServeMux()

	// Extra handlers
	gwmux.HandlePath("POST", "/gateway/v1/song/upload_song", h.UploadSong)
	httpMux.Handle("/", gwmux)
	prefix := os.Getenv("EXPOSED_STORAGE_PREFIX") + "/"
	fileServer := http.FileServer(http.Dir(string(os.Getenv("INTERNAL_STORAGE_PREFIX"))))
	httpMux.Handle(prefix, http.StripPrefix(prefix, fileServer))

	if err := http.Serve(restLis, httpMux); err != nil {
		channel <- fmt.Errorf("could not serve Rest server on port %s: %w", port, err)
	}
}

func (h *Handler) Server() error {
	grpcChannel, restChannel := make(chan error), make(chan error)
	go h.RunGrpcServer(os.Getenv("GRPC_PORT"), grpcChannel)
	go h.RunRestServer(os.Getenv("REST_PORT"), restChannel)

	select {
	case grpcError := <-grpcChannel:
		return grpcError
	case restError := <-restChannel:
		return restError
	}
}
