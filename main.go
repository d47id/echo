package main

import (
	"flag"
	"io/ioutil"

	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	DevelopmentLogger bool `yaml:"DevelopmentLogger"`
}

func main() {
	// Read flags
	configPath := flag.String(
		"config-file",
		"echo.yaml",
		"Path to yaml config file",
	)
	flag.Parse()

	// Read config
	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	cfg := &config{}
	yaml.Unmarshal(data, cfg)

	// Create logger
	var logger *zap.Logger
	if cfg.DevelopmentLogger {
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
		zap.Bool("DevelopmentLogger", cfg.DevelopmentLogger),
	)
}
