package interceptor

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func CheckIPInterceptor(subnet *string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if subnet != nil {
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				if realIP, found := md["ip"]; found {
					ip := net.ParseIP(realIP[0])
					_, subn, err := net.ParseCIDR(*subnet)
					if err != nil {
						return nil, status.Errorf(codes.Unauthenticated, "wrong trusted subnet format: %v", err)
					}
					if !subn.Contains(ip) {
						return nil, status.Errorf(codes.Unauthenticated, "ip not in trusted subnet")
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
