package database

const GenreTableName = "genre"
const SongTableName = "song"
const PlaylistTableName = "playlist"
const PlayListSongTableName = "playlist_song"

type Tabler interface {
	TableName() string
}

func (GenreSchema) TableName() string {
	return GenreTableName
}

func (SongSchema) TableName() string {
	return SongTableName
}

func (PlaylistSchema) TableName() string {
	return PlaylistTableName
}

func (PlaylistSchemaCreate) TableName() string {
	return PlaylistTableName
}

func (PlaylistSchemaPut) TableName() string {
	return PlaylistTableName
}

func (DeletePlaylistSchema) TableName() string {
	return PlaylistTableName
}

func (PlaylistSongSchema) TableName() string {
	return PlayListSongTableName
}

func (PlaylistSongSchemaCreate) TableName() string {
	return PlayListSongTableName
}
