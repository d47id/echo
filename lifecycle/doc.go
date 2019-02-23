// Package lifecycle provides common functionality related to management of the
// life of Go services in a Kubernetes environment. Some might call lifecycle an
// abomination, an affront to the gods of never using "utility" packages. These
// people may be right. TBH, adding the promhttp route option felt a little
// dirty.
//
// However, I find this so tremendously useful, and the individual pieces
// connected enough by the idea of lifecycle management that I'm willing to
// accept its potentially flawed conception.
//
// For a working example, see github.com/d47id/echo
//
// Usage
//
// Create a *lifecycle.Manager by calling lifecycle.New() in func main(). If
// your application wants to support hot config reload, construct a function to
// perform the necessary steps and pass it in as the reloadFunc parameter. New
// also expects a logger implementing the lifecycle.Logger interface. It is
// designed to accept a *zap.SugaredLogger (go.uber.org/zap), but it should be
// simple enough to construct a shim around your logger of choice if you prefer
// something else.
//
// Whenever convenient, call the manager's StartServer() method. It expects a
// port on which to expose the http liveness and readiness probes as well as a
// flag to indicate whether it should also mount promhttp.Handler() at /metrics.
// It returns a function that will shut the server down gracefully. This should
// be pushed onto the defer stack. The endpoints will return 503 Server
// Unavailable until the Healthy() and Ready() methods (respectively) are
// invoked.
//
// The server also exposes version information at /version. The intention is to
// override the values of Version, BuildTime, Branch, and Commit by passing -X
// arguments to -ldflags when the application is built:
//   go build -ldflags "-s -w -X lifecycle.Version=${HASH} -X main.BuildTime=${BUILD_TIME}"
// Any values not overridden will return as "unset".
//
// As you build your application's dependencies in func main(), push functions
// onto the defer stack that will perform any operations necessary to clean the
// dependencies up. Launch any long-running processes in background goroutines
// and push functions onto the defer stack that will gracefully stop them. Once
// your application is ready to start serving requests, call the manager's
// Ready() and Healthy() methods. This will cause the liveness and readiness
// endpoints to start returning 200 OK. The manager is safe for use in
// concurrent goroutines if you want to implement rate-limiting or
// circuit-breaking functionality. Calling the manager's NotHealthy and NotReady
// methods will cause the liveness and readiness probes to return 503 Server
// Unavailable.
//
// The last line of func main() should be a call to the manager's
// WaitForSignals() method. If the application receives SIGINT or SIGTERM, the
// method will invoke NotReady and NotHealthy and return, allowing the cleanup
// operations on the defer stack to execute. If SIGHUP is recieved and a
// reloadFunc was passed to New() when creating the manager, the method will
// invoke the func and continue to block waiting for signals.
package lifecycle
