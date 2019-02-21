package lifecycle

import (
	"context"
	"net/http"
	"time"
)

// StartServer starts an http server exposing liveness, readiness, and
// version info endpoints listening on the given address.
func (m *Manager) StartServer(address string) func() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", m.healthz)
	mux.HandleFunc("/readyz", m.readyz)
	mux.HandleFunc("/version", versionHandler)

	srv := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func(l Logger, srv *http.Server) {
		err := srv.ListenAndServe()
		if err != nil {
			l.Errorw(err.Error())
		}
	}(m.l, srv)

	m.l.Infow("http server started", "address", address)

	return getCancelFunc(m.l, srv)
}

func getCancelFunc(l Logger, srv *http.Server) func() {
	return func() {
		l.Infow("lifecycle http server shutting down...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			l.Errorw(err.Error())
		}
		l.Infow("lifecycle http server stopped")
	}
}
