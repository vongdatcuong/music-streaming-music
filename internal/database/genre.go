package database

import (
	"context"
	"fmt"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
)

func (db *Database) GetGenreOptionsList(ctx context.Context) ([]common.NameValueInt32, error) {
	var schemas []GenreSchema

	result := db.GormClient.Order("name").Find(&schemas)

	if result.Error != nil {
		return []common.NameValueInt32{}, fmt.Errorf("could not get the genre options list: %w", result.Error)
	}

	var optionsList []common.NameValueInt32

	for _, model := range schemas {
		optionsList = append(optionsList, convertGenreSchemaToNameValueInt32(model))
	}

	return optionsList, nil
}
