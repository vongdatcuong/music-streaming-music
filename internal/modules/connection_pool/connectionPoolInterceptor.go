package connection_pool

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TODO: 1 Possible Improvement: only inject token if the resource requires it
type ConnectionPoolInterceptor struct {
}

func NewConnectionPoolInterceptor() *ConnectionPoolInterceptor {
	return &ConnectionPoolInterceptor{}
}

func (cpInterceptor *ConnectionPoolInterceptor) UnaryInjectToken() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		authCtxValue := ctx.Value("authorization")
		if authCtxValue != nil {
			newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", authCtxValue.(string))
			return invoker(newCtx, method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
