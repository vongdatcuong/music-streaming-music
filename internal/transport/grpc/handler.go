package grpc

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
	"google.golang.org/grpc"
)

type Handler struct {
	grpcPbV1.UnimplementedSongServiceServer
	grpcPbV1.UnimplementedPlaylistServiceServer
	grpcPbV1.UnimplementedGenreServiceServer
	songService     SongServiceGrpc
	playlistService PlaylistServiceGrpc
	authInterceptor *AuthInterceptor
	genreService    GenreServiceGrpc
}

func NewHandler(songService SongServiceGrpc, playlistService PlaylistServiceGrpc, authInterceptor *AuthInterceptor, genreService GenreServiceGrpc) *Handler {
	h := &Handler{songService: songService, playlistService: playlistService, authInterceptor: authInterceptor, genreService: genreService}

	return h
}

func (h *Handler) RunGrpcServer(port string, channel chan error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		channel <- fmt.Errorf("could not listen on port %s: %w", port, err)
		return
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(h.authInterceptor.GrpcUnary()))
	grpcPbV1.RegisterSongServiceServer(grpcServer, h)
	grpcPbV1.RegisterPlaylistServiceServer(grpcServer, h)
	grpcPbV1.RegisterGenreServiceServer(grpcServer, h)

	if err := grpcServer.Serve(lis); err != nil {
		channel <- fmt.Errorf("could not server Grpc server on port %s: %w", port, err)
	}
}

func (h *Handler) RunRestServer(port string, channel chan error) {
	/*gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true, // Rest Server to return the same fields as protobuf
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	muxCtx, cancelMuxCtx := context.WithCancel(context.Background())
	defer cancelMuxCtx()
	err := grpcPbV1.RegisterSongServiceHandlerServer(muxCtx, gwmux, h)

	if err != nil {
		channel <- fmt.Errorf("Failed to register Song Rest endpoints: %w", err)
		return
	}

	err = grpcPbV1.RegisterPlaylistServiceHandlerServer(muxCtx, gwmux, h)

	if err != nil {
		channel <- fmt.Errorf("Failed to register Playlist Rest endpoints: %w", err)
		return
	}

	restLis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		channel <- fmt.Errorf("could not listen on port %s: %w", port, err)
		return
	}
	httpMux := http.NewServeMux()

	// Extra handlers
	gwmux.HandlePath("POST", "/api/gateway/v1/song/upload_song", h.UploadSong)
	httpMux.Handle("/", gwmux)
	prefix := os.Getenv("EXPOSED_STORAGE_PREFIX") + "/"
	fileServer := http.FileServer(http.Dir(string(os.Getenv("INTERNAL_STORAGE_PREFIX"))))
	httpMux.Handle(prefix, http.StripPrefix(prefix, fileServer))

	if err := http.Serve(restLis, httpMux); err != nil {
		channel <- fmt.Errorf("could not serve Rest server on port %s: %w", port, err)
	}*/

	restLis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		channel <- fmt.Errorf("could not listen on port %s: %w", port, err)
		return
	}
	httpMux := http.NewServeMux()

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/api/gateway/v1/song/upload_song", h.UploadSong).Methods("POST")
	httpMux.Handle("/", muxRouter)

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
