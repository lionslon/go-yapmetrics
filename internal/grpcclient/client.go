package grpcclient

import (
    "context"
    "time"

    pb "github.com/lionslon/go-yapmetrics/api"
    "google.golang.org/grpc"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.MetricsServiceClient
}

func NewClient(addr string) (*Client, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    return &Client{
        conn:   conn,
        client: pb.NewMetricsServiceClient(conn),
    }, nil
}

func (c *Client) SendMetric(metric *pb.Metric) error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    _, err := c.client.SendMetric(ctx, metric)
    return err
}

func (c *Client) GetMetrics() (*pb.MetricsList, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    return c.client.GetMetrics(ctx, &pb.Empty{})
}

func (c *Client) Close() {
    c.conn.Close()
}
