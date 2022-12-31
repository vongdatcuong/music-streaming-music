package connection_pool

import (
	"fmt"

	grpcPb "github.com/vongdatcuong/music-streaming-protos/go/v1"
	"google.golang.org/grpc"
)

type ConnectionPool struct {
	interceptor      *ConnectionPoolInterceptor
	authConnection   *grpc.ClientConn
	UserClient       grpcPb.UserServiceClient
	PermissionClient grpcPb.PermissionServiceClient
}

func NewConnectionPool(interceptor *ConnectionPoolInterceptor, authServiceAddress string) (*ConnectionPool, error) {
	authConnection, err := grpc.Dial(authServiceAddress, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor.UnaryInjectToken()))

	if err != nil {
		return nil, fmt.Errorf("Error while making connection to Authentication service, %v", err)
	}

	return &ConnectionPool{
		interceptor:      interceptor,
		authConnection:   authConnection,
		UserClient:       grpcPb.NewUserServiceClient(authConnection),
		PermissionClient: grpcPb.NewPermissionServiceClient(authConnection),
	}, nil
}

func (cp *ConnectionPool) CloseAll() error {
	err := cp.authConnection.Close()

	if err != nil {
		return fmt.Errorf("could not close connection to Authenticate service: %w", err)
	}

	return nil
}
