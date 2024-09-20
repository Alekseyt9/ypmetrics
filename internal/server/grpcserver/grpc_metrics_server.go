package grpcserver

import (
	"context"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
	pb "github.com/Alekseyt9/ypmetrics/internal/common/proto"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcMetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	Store storage.Storage
	Log   log.Logger
}

func (s *GrpcMetricsServer) SendBatch(ctx context.Context, in *pb.SendBatchRequest) (*pb.SendBatchResponse, error) {
	var response pb.SendBatchResponse

	mdata := items.MetricItems{}
	for _, m := range in.Metrics {
		switch m.Type {
		case pb.SendBatchRequest_GAUGE:
			mdata.Gauges = append(mdata.Gauges, items.GaugeItem{Name: m.Id, Value: m.Value})
		case pb.SendBatchRequest_COUNTER:
			mdata.Counters = append(mdata.Counters, items.CounterItem{Name: m.Id, Value: int64(m.Value)})
		}
	}

	err := s.Store.SetCounters(ctx, mdata.Counters)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error SetCounters %v", err)
	}

	err = s.Store.SetGauges(ctx, mdata.Gauges)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error SetGauges %v", err)
	}

	return &response, nil
}
