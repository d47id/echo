package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/d47id/echo/api"
	"github.com/d47id/echo/impl"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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

	// build gRPC server
	s := newRPCserver(logger)

	// Build api server implementation
	impl := &impl.EchoImpl{Logger: logger}
	if *relayAddr != "" {
		cl, err := getRPCconn(context.Background(), *relayAddr, logger)
		if err != nil {
			panic(err)
		}
		impl.Client = api.NewEchoClient(cl)
	}

	// Create health server
	hlth := health.NewServer()

	// Register gRPC services with server
	api.RegisterEchoServer(s, impl)
	reflection.Register(s)
	grpc_health.RegisterHealthServer(s, hlth)

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

	// Start http server
	cancel := startHTTPServer(fmt.Sprintf(":%d", *httpPort), logger)
	defer cancel()

	//wait for signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	// set health status to not serving
	hlth.Shutdown()
}
