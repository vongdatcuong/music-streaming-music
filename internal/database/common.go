package database

import (
	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
	"gorm.io/gorm"
)

func Paginate(paginationInfo common.PaginationInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := paginationInfo.Offset
		limit := paginationInfo.Limit

		if offset < 0 {
			offset = 0
		}

		if limit <= 0 {
			limit = constants.PAGINATION_INFO_DEFAULT_LIMIT
		}

		return db.Offset(int(offset)).Limit(int(limit))
	}
}
