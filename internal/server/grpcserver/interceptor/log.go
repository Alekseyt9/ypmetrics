package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/Alekseyt9/ypmetrics/internal/server/log"
)

func LogInterceptor(log log.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		st, _ := status.FromError(err)
		code := st.Code()

		var reqSize int
		if protoMsg, ok := req.(proto.Message); ok {
			reqSize = proto.Size(protoMsg)
		}

		log.Info("gRPC query",
			"method", info.FullMethod,
			"duration", duration,
			"status", code,
			"response_size", reqSize,
		)

		return resp, err
	}
}
