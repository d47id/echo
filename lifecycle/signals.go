package lifecycle

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForSignals blocks until the process recieves SIGINT, SIGTERM, or
// SIGHUP. If SIGINT or SIGTERM is received, the method will return.
// If SIGHUP is received and a ReloadFunc was passed to New(),
// the ReloadFunc will be invoked.
func (m *Manager) WaitForSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
wait:
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP:
			if m.reloadFunc != nil {
				m.reloadFunc()
			}
		default:
			break wait
		}
	}

	// return unavailable from liveness and readiness endpoints
	m.NotReady()
	m.NotHealthy()
}
