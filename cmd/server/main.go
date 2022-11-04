package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-music/internal/database"
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

	/*if err = migration.Down(); err != nil {
		return fmt.Errorf("could not down the migration: %w", err)
	}*/

	grpcHandler := grpcTransport.NewHandler()

	if err := grpcHandler.Server(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Errorln(err)
	}
}
