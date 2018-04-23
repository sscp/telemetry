package collector

import (
	"context"
	"net"
	"strconv"

	pb "github.com/sscp/telemetry/collector/serviceproto"

	"google.golang.org/grpc"
)

// CollectorClient calls and manipulates the CollectorService, used in the call
// command in the CLI
type CollectorClient struct {
	client     pb.CollectorServiceClient
	connection *grpc.ClientConn
}

// CollectorClientConfig is a struct that corresponds the client config in the
// telemetry config file. It specifies which CollectorService to connect to on
// the network
type CollectorClientConfig struct {
	Hostname string
	Port     int
}

// NewCollectorClient constructs a CollectorClient that connects to the
// CollectorService specified in CollectorClientConfig
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

// StartCollector instructs the CollectorService to begin collecting packets
// under the given runName
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

// StopCollector instructs teh CollectorService to stop collecting packets and
// close off current run
func (cc *CollectorClient) StopCollector() (*pb.CollectorStatus, error) {
	req := &pb.StopRequest{}
	status, err := cc.client.StopCollecting(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// GetCollectorStatus returns the status of the CollectorService that includes
// whether the CollectorService is collecting or not and if it is how many
// packets it has collected on the current run
func (cc *CollectorClient) GetCollectorStatus() (*pb.CollectorStatus, error) {
	req := &pb.StatusRequest{}
	status, err := cc.client.GetCollectorStatus(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// Close closes the connection to the CollectorService
func (cc *CollectorClient) Close() {
	cc.connection.Close()
}
