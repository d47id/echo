package impl

import (
	"context"
	"errors"

	"github.com/d47id/echo/api"
)

// EchoServer is the server API for Echo service.
// type EchoServer interface {
// 	Shout(context.Context, *ShoutRequest) (*ShoutReply, error)
// }

// EchoImpl implements api.EchoServer
type EchoImpl struct{}

// Shout implements api.EchoServer.Shout
func (i *EchoImpl) Shout(context.Context, *api.ShoutRequest) (*api.ShoutReply, error) {
	return nil, errors.New("not implemented")
}
