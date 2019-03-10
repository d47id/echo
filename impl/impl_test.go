package impl

import (
	"context"
	"net"
	"testing"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/d47id/echo/api"
)

func TestEchoImpl(t *testing.T) {
	//init logger
	l, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	//initialize bufconn
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)
	defer lis.Close()

	//create grpc server
	srv := &EchoImpl{
		Logger: l,
	}
	s := grpc.NewServer()
	api.RegisterEchoServer(s, srv)
	t.Log("starting test gRPC server")
	go func(t *testing.T, s *grpc.Server, l net.Listener) {
		err := s.Serve(l)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("test gRPC server stopped")
	}(t, s, lis)
	defer s.GracefulStop()

	//create grpc client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithDialer(getDialer(lis)),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := api.NewEchoClient(conn)

	t.Run("messages match", func(t *testing.T) {
		msg := "To gRPC or not to gRPC; is that your question?"
		reply, err := c.Shout(ctx, &api.ShoutRequest{Message: msg})
		if err != nil {
			t.Fatal(err)
		}
		if reply.Message != msg {
			t.Fatalf("unexpected mismatch. wanted %s got %s\n",
				msg, reply.Message)
		}
	})
}

func getDialer(l *bufconn.Listener) func(string, time.Duration) (net.Conn, error) {
	return func(string, time.Duration) (net.Conn, error) {
		return l.Dial()
	}
}
