# music-streaming-music

Music Streaming System - Music Service

- Configure 2 containers: service & db with docker-compose
- Write Database migration scripts
- Create makefile with common commands
- Write protobuf files + Generate pb files
- Configure grpc server
- Implement Database Methods related to Song Table:
  - Get song list with filter and pagination; count total number of records with only filter
- Add Validation for CreateSong and PutSong, using "github.com/go-playground/validator/v10"
- Playlist implementation
- Write a Storage module to save and get files from this server
- Research grpc-gateway and use it to enable Rest endpoints alongside with Grpc Endpoint. Understand more about go channel and go routine in the process
- Service static files on the Rest server
- Implement a REST endpoint to upload song. It will return resource id and resource link for FE to pass into Create/Put Song requests
- Research how to calculate audio file duration. Reach the conclusion to let the FE calculate it and pass the duration to the BE. Also update duration field from uint32 to float
- Define Gorm Schemas and configurations. Able to return song list of playlists using gorm Preload (Join in sql)
- Implement UploadPlaylistSongs endpoint, which allows to update song list of a playlist only
- Use jwt, grpc interceptor, http middleware to implement authentication and authorization
- Migrate Authentication operation to API Gateway. One issue was to defined error response for all Grpc endpoints
- Inject jwt token to grpc interceptor context so that grpc methods have access to jwt token
- Implement Get genre options list endpoint: define protobuf -> Add genre grpc method -> Add genre service -> Add genre database method
- Integrate Cloudinary to storage audio file

Improvements can make

- If fail to add songs to playlist, reverse creating the playlist
