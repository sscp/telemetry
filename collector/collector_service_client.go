package collector

import (
	"context"
	"net"
	"strconv"

	pb "github.com/sscp/telemetry/collector/serviceproto"

	"google.golang.org/grpc"
)

type CollectorClient struct {
	client     pb.CollectorServiceClient
	connection *grpc.ClientConn
}

type CollectorClientConfig struct {
	Hostname string
	Port     int
}

func NewCollectorClient(cfg CollectorClientConfig) (*CollectorClient, error) {
	addr := net.JoinHostPort(cfg.Hostname, strconv.FormatInt(int64(cfg.Port), 10))
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewCollectorServiceClient(conn)
	return &CollectorClient{
		client:     client,
		connection: conn,
	}, nil
}

func (cc *CollectorClient) StartCollector(runName string) (*pb.CollectorStatus, error) {
	req := &pb.StartRequest{
		RunName: runName,
	}
	status, err := cc.client.StartCollecting(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (cc *CollectorClient) StopCollector() (*pb.CollectorStatus, error) {
	req := &pb.StopRequest{}
	status, err := cc.client.StopCollecting(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (cc *CollectorClient) GetCollectorStatus() (*pb.CollectorStatus, error) {
	req := &pb.StatusRequest{}
	status, err := cc.client.GetCollectorStatus(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (cc *CollectorClient) Close() {
	cc.connection.Close()
}
