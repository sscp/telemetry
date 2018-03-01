package collector

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/sscp/telemetry/collector/serviceproto"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:generate protoc -I serviceproto collector_service.proto --go_out=plugins=grpc:serviceproto

const nilPort = 0

type CollectorService struct {
	collector     *Collector
	collectorPort int32
}

type CollectorServiceConfig struct {
	Port      int32
	Collector CollectorConfig
}

func NewCollectorService(cfg CollectorServiceConfig) (*CollectorService, error) {
	collector, err := NewUDPCollector(cfg.Collector)
	if err != nil {
		return nil, err
	}
	return &CollectorService{
		collector:     collector,
		collectorPort: int32(cfg.Collector.Port),
	}, nil
}

func (cs *CollectorService) getStatus() *pb.CollectorStatus {
	status := cs.collector.GetStatus()
	if status.Collecting {
		return &pb.CollectorStatus{
			Collecting:      status.Collecting,
			RunName:         status.RunName,
			PacketsRecorded: status.PacketsProcessed,
			Port:            cs.collectorPort,
		}
	} else {
		return &pb.CollectorStatus{
			Collecting: false,
		}

	}
}

func (cs *CollectorService) StartCollecting(ctx context.Context, req *pb.StartRequest) (*pb.CollectorStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "collectorService/StartCollecting")
	defer span.Finish()
	log.Printf("Starting collection: %v", req)

	// We start a new run with new name if the collector is currently collecting
	if cs.collector.GetStatus().Collecting {
		cs.collector.Close(ctx)
	}

	cs.collector.RecordRun(ctx, req.GetRunName())
	return cs.getStatus(), nil
}

func (cs *CollectorService) StopCollecting(ctx context.Context, req *pb.StopRequest) (*pb.CollectorStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "collectorService/StopCollecting")
	defer span.Finish()

	log.Print("Stopping collection")
	// Only call close if the collector is collecting
	if cs.collector.GetStatus().Collecting {
		cs.collector.Close(ctx)
	}

	return cs.getStatus(), nil
}

func (cs *CollectorService) GetCollectorStatus(ctx context.Context, req *pb.StatusRequest) (*pb.CollectorStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "collectorService/GetCollectorStatus")
	defer span.Finish()
	return cs.getStatus(), nil

}

func RunCollectionService(cfg CollectorServiceConfig) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	srv, err := NewCollectorService(cfg)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}
	pb.RegisterCollectorServiceServer(grpcServer, srv)
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

}
