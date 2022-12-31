package grpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/vongdatcuong/music-streaming-music/internal/modules/connection_pool"
	"github.com/vongdatcuong/music-streaming-music/internal/modules/jwtAuth"
	common_utils "github.com/vongdatcuong/music-streaming-music/internal/modules/utils/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	jwtService            *jwtAuth.JwtService
	accessiblePermissions map[string][]string
	connectionPool        *connection_pool.ConnectionPool
}

func NewAuthInterceptor(jwtService *jwtAuth.JwtService, connectionPool *connection_pool.ConnectionPool) *AuthInterceptor {
	return &AuthInterceptor{jwtService: jwtService, connectionPool: connectionPool}
}

func (interceptor *AuthInterceptor) GrpcUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		injectedCtx, err, _ := interceptor.authorize(ctx, md["authorization"], info.FullMethod, EndPointPermissions, EndPointNoAuthentication)

		if err != nil {
			return getRespective403Response(info.FullMethod), nil
		}

		return handler(injectedCtx, req)
	}
}

func (interceptor *AuthInterceptor) HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err, errCode := interceptor.authorize(r.Context(), r.Header["Authorization"], r.URL.Path, HttpEndPointPermissions, HttpEndPointNoAuthentication)

		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, errCode, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, authHeader []string, path string, permissionsMap map[string][]string, noAuthenMap map[string]bool) (context.Context, error, uint32) {
	if noAuthenMap[path] {
		return ctx, nil, 0
	}

	accessToken, err := parseAuthorizationHeader(authHeader)

	if err != nil {
		return ctx, err, 1
	}

	claims, err := interceptor.jwtService.ValidateToken(accessToken)

	if err != nil {
		return ctx, err, 1
	}

	requiredPerm := permissionsMap[path]
	var firstRequiredPermName string

	if requiredPerm != nil && len(requiredPerm) > 0 {
		firstRequiredPermName = requiredPerm[0]
	}

	newCtx := interceptor.jwtService.InjectToken(ctx, accessToken)
	res, err := interceptor.connectionPool.PermissionClient.CheckUserPermission(newCtx, &grpcPbV1.CheckUserPermissionRequest{UserId: claims.UserID, PermissionName: firstRequiredPermName})

	if err != nil {
		return newCtx, err, 1
	}

	if res == nil || res.HasPermission == nil || !*res.HasPermission {
		return newCtx, fmt.Errorf("you have no permission to access this resource"), 403
	}

	return newCtx, nil, 0
}

func parseAuthorizationHeader(values []string) (string, error) {
	if values == nil || len(values) == 0 {
		return "", fmt.Errorf("invalid authorization header")
	}
	authHeader := values[0]
	authHeaderParts := strings.Split(authHeader, " ")

	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header")
	}

	return authHeaderParts[1], nil
}

// TODO: improve this
func getRespective403Response(path string) any {
	errCode, errMsg := common_utils.GetUInt32Pointer(403), common_utils.GetStringPointer("You have no permission to access this resource")

	if path == songServicePath+"/GetSongList" {
		return &grpcPbV1.GetSongListResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == songServicePath+"/GetSongDetails" {
		return &grpcPbV1.GetSongDetailsResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == songServicePath+"/CreateSong" {
		return &grpcPbV1.CreateSongResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == songServicePath+"/PutSong" {
		return &grpcPbV1.PutSongResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == songServicePath+"/DeleteSong" {
		return &grpcPbV1.DeleteSongResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == playlistServicePath+"/GetPlaylistList" {
		return &grpcPbV1.GetPlaylistListResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == playlistServicePath+"/GetPlaylistDetails" {
		return &grpcPbV1.GetPlaylistDetailsResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == playlistServicePath+"/CreatePlaylist" {
		return &grpcPbV1.CreatePlaylistResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == playlistServicePath+"/PutPlaylist" {
		return &grpcPbV1.PutPlaylistResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == playlistServicePath+"/DeletePlaylist" {
		return &grpcPbV1.DeletePlaylistResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == playlistServicePath+"/UpdatePlaylistSongs" {
		return &grpcPbV1.UpdatePlaylistSongsResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	}
	return nil
}
