package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/d47id/echo/lifecycle"

	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
}

func main() {
	// Read flags
	configPath := flag.String(
		"config-file",
		"echo.yaml",
		"Path to yaml config file",
	)
	httpPort := flag.Int(
		"http-port",
		3000,
		"Port on which to expose liveness, readiness, and version info endpoints",
	)
	devLogger := flag.Bool(
		"dev-logger",
		false,
		"Use development logger",
	)
	flag.Parse()

	// Load config
	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	cfg := &config{}
	yaml.Unmarshal(data, cfg)

	// Create logger
	var logger *zap.Logger
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
		zap.Int("http-port", *httpPort),
		zap.Bool("dev-logger", *devLogger),
	)

	// Create lifecycle manager
	mgr := lifecycle.New(logger.Sugar(), nil)

	// Start http server
	cancel := mgr.StartServer(fmt.Sprintf(":%d", *httpPort))
	defer cancel()

	// Start gRPC server
	// NOT IMPLEMENTED

	// Set ready, healthy
	mgr.Healthy()
	mgr.Ready()

	//wait for signals
	mgr.WaitForSignals()
}
