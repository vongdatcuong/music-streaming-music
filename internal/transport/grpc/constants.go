package grpc

// Songs
const songPermPrefix = "music_streaming.song"

var songPermissions map[string]string = map[string]string{
	"READ":  songPermPrefix + ".read",
	"WRITE": songPermPrefix + ".write",
}

// Playlists
const playlistPermPrefix = "music_streaming.playlist"

var playlistPermissions map[string]string = map[string]string{
	"READ":  playlistPermPrefix + ".read",
	"WRITE": playlistPermPrefix + ".write",
}

// Genre
const genrePermPrefix = "music_streaming.genre"

var genrePermissions map[string]string = map[string]string{
	"READ":  genrePermPrefix + ".read",
	"WRITE": genrePermPrefix + ".write",
}

// Endpoints
const songServicePath string = "/music_streaming.music.song.SongService"
const playlistServicePath string = "/music_streaming.music.playlist.PlaylistService"
const genreServicePath string = "/music_streaming.music.genre.GenreService"

var EndPointPermissions map[string][]string = map[string][]string{
	// Song
	songServicePath + "/CreateSong": {songPermissions["WRITE"]},
	songServicePath + "/PutSong":    {songPermissions["WRITE"]},
	songServicePath + "/DeleteSong": {songPermissions["WRITE"]},
	// Playlist
	playlistServicePath + "/GetPlaylistList":     {playlistPermissions["READ"]},
	playlistServicePath + "/GetPlaylistDetails":  {playlistPermissions["READ"]},
	playlistServicePath + "/CreatePlaylist":      {playlistPermissions["WRITE"]},
	playlistServicePath + "/PutPlaylist":         {playlistPermissions["WRITE"]},
	playlistServicePath + "/DeletePlaylist":      {playlistPermissions["WRITE"]},
	playlistServicePath + "/UpdatePlaylistSongs": {playlistPermissions["WRITE"]},

	// TODO: Maybe can have a separate set of endpoint for user client (non-admin) to call, which won't check for permission possessed by user
}
var EndPointNoAuthentication map[string]bool = map[string]bool{}

// Http
/*const httpPath = "/api/gateway/v1"
const httpSongPath = httpPath + "/song"
const httpPlaylistPath = httpPath + "/playlist"

var HttpEndPointPermissions map[string][]string = map[string][]string{
	// Song
	httpSongPath + "/list":        {songPermissions["READ"]},
	httpSongPath + "/details":     {songPermissions["READ"]},
	httpSongPath + "/create_song": {songPermissions["WRITE"]},
	httpSongPath + "/put_song":    {songPermissions["WRITE"]},
	httpSongPath + "/delete_song": {songPermissions["WRITE"]},

	// Playlist
	httpPlaylistPath + "/list":                  {playlistPermissions["READ"]},
	httpPlaylistPath + "/details":               {playlistPermissions["READ"]},
	httpPlaylistPath + "/create_playlist":       {playlistPermissions["WRITE"]},
	httpPlaylistPath + "/put_playlist":          {playlistPermissions["WRITE"]},
	httpPlaylistPath + "/delete_playlist":       {playlistPermissions["WRITE"]},
	httpPlaylistPath + "/update_playlist_songs": {playlistPermissions["WRITE"]},
}
var HttpEndPointNoAuthentication map[string]bool = map[string]bool{}
*/
