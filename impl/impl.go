package impl

import (
	"context"

	"github.com/d47id/echo/api"
)

// EchoImpl implements api.EchoServer
type EchoImpl struct{}

// Shout implements api.EchoServer.Shout
func (i *EchoImpl) Shout(ctx context.Context, req *api.ShoutRequest) (*api.ShoutReply, error) {
	return &api.ShoutReply{
		Message: req.Message,
	}, nil
}
