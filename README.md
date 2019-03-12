The goal of this project is to serve as a boilerplate starting point for writing Go gRPC servers for Kubernetes. 

Requirements:
* Configuration via file
* Logging
* Graceful shutdown
* Liveness and readiness http endpoints
* Prometheus metrics
* Working gRPC "echo" service
* gRPC reflection
* OpenTracing header propagation

TODO:
* Implement gRPC health check proto
* Use ctx zap logger
* Deprecate lifecycle module