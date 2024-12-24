package grpcserver

import (
    "context"
    "sync"
    pb "github.com/lionslon/go-yapmetrics/api"
    "google.golang.org/grpc"
    "net"
    "log"
)

type Server struct {
    pb.UnimplementedMetricsServiceServer
    metrics sync.Map
}

func NewServer() *Server {
    return &Server{}
}

func (s *Server) SendMetric(ctx context.Context, metric *pb.Metric) (*pb.Empty, error) {
    log.Printf("[gRPC] Received metric: %+v", metric)
    s.metrics.Store(metric.Id, metric)
    return &pb.Empty{}, nil
}

func (s *Server) GetMetrics(ctx context.Context, empty *pb.Empty) (*pb.MetricsList, error) {
    var metrics []*pb.Metric
    s.metrics.Range(func(_, value interface{}) bool {
        metrics = append(metrics, value.(*pb.Metric))
        return true
    })
    return &pb.MetricsList{Metrics: metrics}, nil
}

func StartGRPCServer(addr string, srv *Server) error {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }

    grpcServer := grpc.NewServer()
    pb.RegisterMetricsServiceServer(grpcServer, srv)
    return grpcServer.Serve(listener)
}
