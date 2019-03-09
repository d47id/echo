package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/d47id/echo/api"
	"github.com/d47id/echo/impl"
	"github.com/d47id/lifecycle"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Read flags
	configPath := flag.String(
		"config-file",
		"echo.yaml",
		"Path to yaml config file",
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
		zap.String("config-file", *configPath),
		zap.Int("grpc-port", *grpcPort),
		zap.Int("http-port", *httpPort),
		zap.Bool("dev-logger", *devLogger),
	)

	// Create lifecycle manager
	mgr := lifecycle.New(logger.Sugar(), nil)

	// Start http server
	cancel := mgr.StartServer(fmt.Sprintf(":%d", *httpPort), false)
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

	// Start gRPC server
	s := grpc.NewServer()
	api.RegisterEchoServer(s, &impl.EchoImpl{Logger: logger})
	reflection.Register(s)
	go func(l *zap.Logger, lis net.Listener, s *grpc.Server) {
		err := s.Serve(lis)
		if err != nil {
			l.Error(err.Error())
		}
	}(logger, lis, s)
	logger.Info("gRPC server started")

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
