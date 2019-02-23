package lifecycle

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartServer starts an http server exposing liveness, readiness, and
// version info endpoints listening on the given address. If metrics is
// true, promhttp.Handler() will be mounted at /metrics
func (m *Manager) StartServer(address string, metrics bool) func() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", m.healthz)
	mux.HandleFunc("/readyz", m.readyz)
	mux.HandleFunc("/version", versionHandler)

	if metrics {
		mux.Handle("/metrics", promhttp.Handler())
	}

	srv := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func(l Logger, srv *http.Server) {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			l.Errorw(err.Error())
		}
	}(m.l, srv)

	m.l.Infow(
		"http server started",
		"address",
		address,
		"metrics",
		metrics,
	)

	return getCancelFunc(m.l, srv)
}

func getCancelFunc(l Logger, srv *http.Server) func() {
	return func() {
		err := srv.Shutdown(context.Background())
		if err != nil {
			l.Errorw(err.Error())
		}
		l.Infow("lifecycle http server stopped")
	}
}
