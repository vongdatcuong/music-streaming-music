package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-music/internal/database"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/connection_pool"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/genre"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/jwtAuth"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/playlist"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/storage"
	"github.com/vongdatcuong/music-streaming-music/internal/transport/grpc"
	grpcTransport "github.com/vongdatcuong/music-streaming-music/internal/transport/grpc"
)

func Run() error {
	var err error
	db, err := database.NewDatabase()

	if err != nil {
		return err
	}

	// Ping DB
	if err := db.Client.DB.Ping(); err != nil {
		return fmt.Errorf("could not ping the database: %w", err)
	}

	_, err = db.MigrateDB()

	if err != nil {
		return err
	}

	// Initiate Connection Pool
	cpInterceptor := connection_pool.NewConnectionPoolInterceptor()
	connectionPool, err := connection_pool.NewConnectionPool(cpInterceptor, os.Getenv("AUTHENTICATION_SERVICE_ADDRESS"))

	if err != nil {
		return err
	}

	defer connectionPool.CloseAll()

	storageService, err := storage.NewService()

	if err != nil {
		return fmt.Errorf("could not initiate connection to storage service: %w", err)
	}

	genreService := genre.NewService(db)
	songService := song.NewService(db, storageService)
	playlistService := playlist.NewService(db)
	jwtAuthService := jwtAuth.NewService(os.Getenv("JWT_SECRET_KEY"), 6*time.Hour)
	authInterceptor := grpc.NewAuthInterceptor(jwtAuthService, connectionPool)
	grpcHandler := grpcTransport.NewHandler(songService, playlistService, authInterceptor, genreService)

	if err := grpcHandler.Server(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		logrus.Errorln(err)
	}
}
