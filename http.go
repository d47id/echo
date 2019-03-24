package main

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func startHTTPServer(address string, l *zap.Logger) func() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/version", versionHandler)

	srv := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func(l *zap.Logger, srv *http.Server) {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			l.Error(err.Error())
		}
	}(l, srv)

	l.Info(
		"http server started",
		zap.String("address", address),
	)

	return getHTTPCancelFunc(l, srv)
}

func getHTTPCancelFunc(l *zap.Logger, srv *http.Server) func() {
	return func() {
		err := srv.Shutdown(context.Background())
		if err != nil {
			l.Error(err.Error())
		}
		l.Info("http server stopped")
	}
}
