package interceptor

import (
	"context"

	"github.com/Alekseyt9/ypmetrics/internal/common/hash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func SetHashInterceptor(hashKey *string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if hashKey != nil {
			data, err := proto.Marshal(req.(proto.Message))
			if err != nil {
				return status.Errorf(codes.Internal, "failed to marshal request: %v", err)
			}
			out := hash.HashSHA256(data, []byte(*hashKey))
			ctx = metadata.AppendToOutgoingContext(ctx, "hash", out)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
