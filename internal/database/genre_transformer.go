package database

import "github.com/vongdatcuong/music-streaming-music/internal/modules/common"

func convertGenreSchemaToNameValueInt32(schema GenreSchema) common.NameValueInt32 {
	return common.NameValueInt32{
		Name:  schema.Name,
		Value: int32(schema.GenreID),
	}
}
