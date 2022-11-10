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
