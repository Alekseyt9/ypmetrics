package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/Alekseyt9/ypmetrics/internal/common/hash"
)

func CheckHashInterceptor(hashKey *string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if hashKey != nil {
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				if hashVal, found := md["hash"]; found {
					reqBytes, err := proto.Marshal(req.(proto.Message))
					if err != nil {
						return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
					}
					if hashVal[0] != hash.HashSHA256(reqBytes, []byte(*hashKey)) {
						return nil, status.Errorf(codes.Internal, "hash check error")
					}
				}
			}
		}

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
