package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-music/internal/database"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/playlist"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/storage"
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

	storageService := storage.NewService()
	songService := song.NewService(db, storageService)
	playlistService := playlist.NewService(db)

	grpcHandler := grpcTransport.NewHandler(songService, playlistService)

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
