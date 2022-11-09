package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/song"
	database_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/database"
	time_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/time"
)

// 12 fields
type SongRow struct {
	SongID       uint64
	Name         string
	Genre        common.NameValueInt32
	Artist       string
	Duration     uint32
	Language     constants.LANGUAGE_ENUM
	Rating       float32
	ResourceID   uint64
	ResourceLink string
	CreatedAt    uint64
	UpdatedAt    uint64
	Status       constants.ACTIVE_STATUS
}

func convertSongRowToSong(songRow SongRow) song.Song {
	return song.Song{
		SongID:       songRow.SongID,
		Name:         songRow.Name,
		Genre:        songRow.Genre,
		Artist:       songRow.Artist,
		Duration:     songRow.Duration,
		Language:     songRow.Language,
		Rating:       songRow.Rating,
		ResourceID:   songRow.ResourceID,
		ResourceLink: songRow.ResourceLink,
		CreatedAt:    songRow.CreatedAt,
		UpdatedAt:    songRow.UpdatedAt,
		Status:       songRow.Status,
	}
}

func (db *Database) GetSongList(ctx context.Context, pagination common.PaginationInfo, filter song.SongListFilter) ([]song.Song, uint64, error) {
	equalQueries := make(map[string]any)
	greaterEqualQueries := make(map[string]any)
	lessEqualQueries := make(map[string]any)

	if filter.Name != nil {
		equalQueries["s.name"] = *filter.Name
	}
	if filter.Genre != nil {
		equalQueries["s.genre"] = *filter.Genre
	}
	if filter.Artist != nil {
		equalQueries["s.artist"] = *filter.Artist
	}
	if filter.Duration != nil {
		equalQueries["s.duration"] = *filter.Duration
	}
	if filter.Language != nil {
		equalQueries["s.language"] = *filter.Language
	}

	if filter.CreatedTimeFrom != nil {
		greaterEqualQueries["s.created_at"] = *filter.CreatedTimeFrom
	}

	if filter.CreatedTimeTo != nil {
		lessEqualQueries["s.created_at"] = *filter.CreatedTimeTo
	}

	equalWheres, equalValues := database_utils.BuildSqlQueryStr(equalQueries, " = ")
	greaterEqualWheres, greaterEqualValues := database_utils.BuildSqlQueryStr(greaterEqualQueries, " >= ")
	lessEqualWheres, lessEqualValues := database_utils.BuildSqlQueryStr(lessEqualQueries, " <= ")

	var wheres []string
	wheres = append(equalWheres, greaterEqualWheres...)
	wheres = append(wheres, lessEqualWheres...)

	var values []any
	values = append(equalValues, greaterEqualValues...)
	values = append(values, lessEqualValues...)

	// Count total songs
	var totalCount uint64
	countQueryStr := database_utils.GetSqlWhereClause(strings.Join(wheres, " AND "))
	countSql := `	SELECT COUNT(*)
								FROM Song s ` + countQueryStr
	rows2, err := db.Client.QueryContext(ctx, countSql, values...)

	if err != nil {
		return []song.Song{}, 0, fmt.Errorf("could not count number of songs: %w", err)
	}
	for rows2.Next() {
		if err := rows2.Scan(&totalCount); err != nil {
			return []song.Song{}, 0, fmt.Errorf("could not count number of songs: %w", err)
		}
	}

	// Pagination
	offset := pagination.Offset
	var limit uint64

	if pagination.Limit == 0 {
		limit = constants.PAGINATION_INFO_DEFAULT_LIMIT
	} else {
		limit = pagination.Limit
	}
	values = append(values, offset)
	values = append(values, limit)

	queryStr := database_utils.GetSqlWhereClause(strings.Join(wheres, " AND "))

	sql := `SELECT s.*, g.name
					FROM Song s 
					INNER JOIN Genre g ON s.genre = g.genre_id ` + queryStr + ` LIMIT ?, ?`
	rows, err := db.Client.QueryContext(ctx, sql, values...)

	if err != nil {
		return []song.Song{}, 0, fmt.Errorf("could not get song list: %w", err)
	}
	defer rows.Close()

	var songs []song.Song
	for rows.Next() {
		songRow := song.Song{}
		err := rows.Scan(&songRow.SongID, &songRow.Name, &songRow.Genre.Value, &songRow.Artist, &songRow.Duration, &songRow.Language, &songRow.Rating,
			&songRow.ResourceID, &songRow.ResourceLink, &songRow.CreatedAt, &songRow.UpdatedAt, &songRow.Status, &songRow.Genre.Name)

		if err != nil {
			return []song.Song{}, 0, fmt.Errorf("could not scan song row: %w", err)
		}

		songs = append(songs, convertSongRowToSong(SongRow(songRow)))
	}

	return songs, totalCount, nil
}

func (db *Database) GetSongDetails(ctx context.Context, id uint64) (song.Song, error) {
	var songRow SongRow
	sql := `SELECT s.*, g.name 
					FROM Song s 
					INNER JOIN Genre g ON s.genre = g.genre_id
	 				WHERE s.song_id = ? AND s.status = ?`
	row := db.Client.QueryRowContext(ctx, sql, id, constants.ACTIVE_STATUS_ACTIVE)
	err := row.Scan(&songRow.SongID, &songRow.Name, &songRow.Genre.Value, &songRow.Artist, &songRow.Duration, &songRow.Language, &songRow.Rating,
		&songRow.ResourceID, &songRow.ResourceLink, &songRow.CreatedAt, &songRow.UpdatedAt, &songRow.Status, &songRow.Genre.Name)

	if err != nil {
		return song.Song{}, fmt.Errorf("could not fetch the song id %d: %w", id, err)
	}
	return convertSongRowToSong(songRow), nil
}

func (db *Database) CreateSong(ctx context.Context, newSong song.Song) (song.Song, error) {
	songRow := SongRow{
		Name:         newSong.Name,
		Genre:        newSong.Genre,
		Artist:       newSong.Artist,
		Duration:     newSong.Duration,
		Language:     newSong.Language,
		Rating:       newSong.Rating,
		ResourceID:   1,
		ResourceLink: "empty",
		CreatedAt:    time_utils.GetCurrentUnixTime(),
		UpdatedAt:    time_utils.GetCurrentUnixTime(),
		Status:       constants.ACTIVE_STATUS_ACTIVE,
	}
	sql := `INSERT INTO Song(name, genre, artist, duration, language, rating, resource_id, resource_link,
						created_at, updated_at, status) VALUES (:name, :genre.value, :artist, :duration, :language, :rating, :resourceid, :resourcelink, :createdat,
						:updatedat, :status)`

	row, err := db.Client.NamedExecContext(ctx, sql, songRow)

	if err != nil {
		return song.Song{}, fmt.Errorf("could not insert new song: %w", err)
	}
	lastInsertedID, err := row.LastInsertId()

	if err != nil {
		return song.Song{}, fmt.Errorf("could not get the last insert ID: %w", err)
	}

	// Update song param with new song ID
	newSong.SongID = uint64(lastInsertedID)

	return newSong, nil
}

// Put Song
func (db *Database) PutSong(ctx context.Context, existingSong song.Song) (song.Song, error) {
	songRow := SongRow{
		SongID:       existingSong.SongID,
		Name:         existingSong.Name,
		Genre:        existingSong.Genre,
		Artist:       existingSong.Artist,
		Duration:     existingSong.Duration,
		Language:     existingSong.Language,
		Rating:       existingSong.Rating,
		ResourceID:   existingSong.ResourceID,
		ResourceLink: existingSong.ResourceLink,
		CreatedAt:    existingSong.CreatedAt,
		UpdatedAt:    time_utils.GetCurrentUnixTime(),
		Status:       existingSong.Status,
	}
	sql := `
				UPDATE Song 
				SET name = :name, genre = :genre.value, artist = :artist, duration = :duration, language = :language, rating = :rating, 
						resource_id = :resourceid, resource_link = :resourcelink, created_at = :createdat, updated_at = :updatedat, status = :status
				WHERE song_id = :songid`

	_, err := db.Client.NamedExecContext(ctx, sql, songRow)

	if err != nil {
		return song.Song{}, fmt.Errorf("could not update song: %w", err)
	}

	return existingSong, nil
}

func (db *Database) DeleteSong(ctx context.Context, id uint64) error {
	sql := "DELETE FROM Song WHERE song_id = ?"
	_, err := db.Client.ExecContext(ctx, sql, id)

	if err != nil {
		return fmt.Errorf("could not delete song with id %d: %w", id, err)
	}

	return nil
}
