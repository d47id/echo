package impl

import (
	"context"

	"github.com/d47id/echo/api"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

// EchoImpl implements api.EchoServer
type EchoImpl struct {
	Logger *zap.Logger
}

// Shout implements api.EchoServer.Shout
func (i *EchoImpl) Shout(ctx context.Context, req *api.ShoutRequest) (*api.ShoutReply, error) {
	i.Logger.Info("request", zap.String("message", req.Message))
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			i.Logger.Info("metadata", zap.Strings(k, v))
		}
	}
	return &api.ShoutReply{
		Message: req.Message,
	}, nil
}
