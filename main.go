package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/d47id/echo/api"
	"github.com/d47id/echo/impl"
	"github.com/d47id/lifecycle"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
    "github.com/opentracing/opentracing-go"
)

func main() {
	// Read flags
	relayAddr := flag.String(
		"relay-addr",
		"",
		"Relay request through downstream echo server",
	)
	grpcPort := flag.Int(
		"grpc-port",
		3000,
		"Port on which to expose gRPC services",
	)
	httpPort := flag.Int(
		"http-port",
		4000,
		"Port on which to expose liveness, readiness, and version info endpoints",
	)
	devLogger := flag.Bool(
		"dev-logger",
		false,
		"Use development logger",
	)
	flag.Parse()

	// Create logger
	var logger *zap.Logger
	var err error
	if *devLogger {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info(
		"Starting up...",
		zap.Int("grpc-port", *grpcPort),
		zap.Int("http-port", *httpPort),
		zap.Bool("dev-logger", *devLogger),
	)

	// Create lifecycle manager
	mgr := lifecycle.New(logger.Sugar(), nil)

	// Start http server
	cancel := mgr.StartServer(fmt.Sprintf(":%d", *httpPort), true)
	defer cancel()

	// Start tcp listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		panic(err)
	}
	logger.Info(
		"tcp listener open",
		zap.String("address", lis.Addr().String()),
	)

	defer func(l *zap.Logger, lis net.Listener) {
		lis.Close()
		l.Info("tcp listener closed")
	}(logger, lis)

	// Configure middleware
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_zap.ReplaceGrpcLogger(logger)
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	
	// Build api server implementation
	impl := &impl.EchoImpl{Logger: logger}
	if *relayAddr != "" {
		cl, err := getConn(context.Background(), *relayAddr, logger)
		if err != nil {
			panic(err)
		}
		impl.Client = api.NewEchoClient(cl)
	}

	// Register gRPC services with server
	api.RegisterEchoServer(s, impl)
	reflection.Register(s)
	grpc_prometheus.Register(s)

	// Start gRPC server
	go func(l *zap.Logger, lis net.Listener, s *grpc.Server) {
		err := s.Serve(lis)
		if err != nil {
			l.Error(err.Error())
		}
	}(logger, lis, s)
	logger.Info("gRPC server started")

	// Stop server on process exit
	defer func(l *zap.Logger, s *grpc.Server) {
		s.GracefulStop()
		l.Info("gRPC server stopped")
	}(logger, s)

	// Set ready, healthy
	mgr.Healthy()
	mgr.Ready()

	//wait for signals
	mgr.WaitForSignals()
}

func getConn(ctx context.Context, addr string, 
	l *zap.Logger) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_prometheus.StreamClientInterceptor,
			grpc_zap.StreamClientInterceptor(l),
			grpc_opentracing.StreamClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_prometheus.UnaryClientInterceptor,
			grpc_zap.UnaryClientInterceptor(l),
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		)),
	)
}
