package collector

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/sscp/telemetry/collector/serviceproto"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
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
	Port       int32
	Collector  CollectorConfig
	JaegerAddr string
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
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

	// Setup tracing if enabled
	if cfg.JaegerAddr != "" {
		log.Printf("Tracing enabled. Sending spans to %v", cfg.JaegerAddr)
		// Sample configuration for testing. Use constant sampling to sample every trace
		jaegerCfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LocalAgentHostPort: cfg.JaegerAddr,
			},
		}

		// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
		// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
		// frameworks.
		jLogger := jaegerlog.StdLogger
		jMetricsFactory := metrics.NullFactory

		// Initialize tracer with a logger and a metrics factory
		closer, err := jaegerCfg.InitGlobalTracer(
			"telemetryCollectorService",
			jaegercfg.Logger(jLogger),
			jaegercfg.Metrics(jMetricsFactory),
		)
		if err != nil {
			log.Printf("Could not initialize jaeger tracer: %s", err.Error())
			return
		}
		defer closer.Close()

	}

	grpcServer.Serve(lis)

}
