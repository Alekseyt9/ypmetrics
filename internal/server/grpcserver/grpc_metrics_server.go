package grpcserver

import (
	"context"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
	pb "github.com/Alekseyt9/ypmetrics/internal/common/proto"
	"github.com/Alekseyt9/ypmetrics/internal/server/log"
	"github.com/Alekseyt9/ypmetrics/internal/server/storage"
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

	_, err := s.Store.GetCounters(ctx)
	if err != nil {
		//http.Error(w, "error GetCounters", http.StatusBadRequest)
	}

	_, err = s.Store.GetGauges(ctx)
	if err != nil {
		//http.Error(w, "error GetGauges", http.StatusBadRequest)
	}

	return &response, nil
}
