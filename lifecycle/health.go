package lifecycle

import (
	"net/http"
)

// Ready sets the readiness endpoint to return OK
func (m *Manager) Ready() {
	m.rmx.Lock()
	m.ready = true
	m.rmx.Unlock()
}

// NotReady sets the readiness endpoint to return Server Unavailable
func (m *Manager) NotReady() {
	m.rmx.Lock()
	m.ready = false
	m.rmx.Unlock()
}

func (m *Manager) readyz(w http.ResponseWriter, _ *http.Request) {
	if m.ready {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w,
		http.StatusText(http.StatusServiceUnavailable),
		http.StatusServiceUnavailable,
	)
}

// Healthy sets the liveness endpoint to return OK
func (m *Manager) Healthy() {
	m.hmx.Lock()
	m.healthy = true
	m.hmx.Unlock()
}

// NotHealthy sets the liveness endpoint to return Server Unavailable
func (m *Manager) NotHealthy() {
	m.hmx.Lock()
	m.healthy = false
	m.hmx.Unlock()
}

func (m *Manager) healthz(w http.ResponseWriter, _ *http.Request) {
	if m.healthy {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w,
		http.StatusText(http.StatusServiceUnavailable),
		http.StatusServiceUnavailable,
	)
}
