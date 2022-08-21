package internalgrpc

import (
	context "context"
	"time"

	grpc "google.golang.org/grpc"
)

func loggingMiddleware(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		response, err := handler(ctx, req)
		logger.LogGRPCRequest(req, info.FullMethod, time.Since(start))

		return response, err
	}
}
