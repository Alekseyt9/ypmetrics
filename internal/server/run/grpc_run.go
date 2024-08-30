package run

import (
	"net"

	"github.com/Alekseyt9/ypmetrics/internal/server/config"
	"github.com/Alekseyt9/ypmetrics/internal/server/grpcserver"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"google.golang.org/grpc"

	pb "github.com/Alekseyt9/ypmetrics/internal/common/proto"
)

func grpcServerStart(store storage.Storage, cfg *config.Config, log log.Logger, stop chan error) {
	go func() {
		listen, err := net.Listen("tcp", *cfg.GRPCAddress)
		if err != nil {
			stop <- err
			return
		}

		s := grpc.NewServer()
		pb.RegisterMetricsServiceServer(s, &grpcserver.GrpcMetricsServer{Store: store, Log: log})

		if err := s.Serve(listen); err != nil {
			stop <- err
			return
		}
	}()
}
