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

	ctx := context.Background()
	a, cancelA := startTestServer(ctx, l, nil, t)
	defer cancelA()
	b, cancelB := startTestServer(ctx, l, a, t)
	defer cancelB()

	t.Run("messages match", func(t *testing.T) {
		msg := "To gRPC or not to gRPC; is that your question?"
		reply, err := a.Shout(ctx, &api.ShoutRequest{Message: msg})
		if err != nil {
			t.Fatal(err)
		}
		if reply.Message != msg {
			t.Fatalf("unexpected mismatch. wanted %s got %s\n",
				msg, reply.Message)
		}
	})

	t.Run("relayed messages match", func(t *testing.T) {
		msg := "To gRPC or not to gRPC; is that your question?"
		reply, err := b.Shout(ctx, &api.ShoutRequest{Message: msg})
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

func startTestServer(
	ctx context.Context,
	l *zap.Logger,
	relayClient api.EchoClient,
	t *testing.T) (
	client api.EchoClient,
	cancel func()) {

	//initialize bufconn
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	//create grpc server
	srv := &EchoImpl{
		Logger: l,
		Client: relayClient,
	}
	s := grpc.NewServer()
	api.RegisterEchoServer(s, srv)
	t.Log("starting test gRPC server")
	go func(t *testing.T, s *grpc.Server, l net.Listener) {
		err := s.Serve(l)
		if err != nil {
			t.Fatal(err)
		}
	}(t, s, lis)

	//create grpc client
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithDialer(getDialer(lis)),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatal(err)
	}

	return api.NewEchoClient(conn), getCancelFunc(lis, s, conn, t)
}

func getCancelFunc(lis net.Listener, s *grpc.Server,
	conn *grpc.ClientConn, t *testing.T) func() {
	return func() {
		err := conn.Close()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("client connection closed")
		s.GracefulStop()
		t.Log("test gRPC server stopped")
		err = lis.Close()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("bufconn listener closed")
	}
}
