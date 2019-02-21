package lifecycle

import "sync"

// Manager encapsulates the functionality provided by the liveness package.
// It presents liveness and readiness endpoints with a mechanism to declare
// the system Ready or NotReady and Healthy or NotHealthy. It also contains
// a mechanism to block waiting for OS signals and optionally reload config.
type Manager struct {
	ready        bool
	rmx          *sync.Mutex
	healthy      bool
	hmx          *sync.Mutex
	configReload func()
	l            Logger
}

// Logger is the interface the lifecycle package uses to log. It was designed
// to accept a *zap.SugaredLogger, but should be easily adaptible to other
// packages.
type Logger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

// New returns a new Manager ready to be used. Logger is not optional
// and must be provided. ReloadFunc is optional, and if provided, will
// be invoked whenever SIGHUP is trapped after invoking WaitForSignals()
func New(logger Logger, ReloadFunc func()) *Manager {
	return &Manager{
		rmx:          &sync.Mutex{},
		hmx:          &sync.Mutex{},
		configReload: ReloadFunc,
		l:            logger,
	}
}
