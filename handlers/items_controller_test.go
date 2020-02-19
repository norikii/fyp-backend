package handlers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tatrasoft/fyp-backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024
var listener *bufconn.Listener
var ctx context.Context

func init() {
	ctx = context.Background()
	listener = bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	proto.RegisterItemServiceServer(srv, &ItemsServerService{})

	go func() {
		if err := srv.Serve(listener); err !=nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func buffDialer(context.Context, string) (net.Conn, error) {
	return listener.Dial()
}


func TestItemsServerService_CreateItem(t *testing.T) {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(buffDialer) ,grpc.WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, conn)

	defer conn.Close()
	client := proto.NewItemServiceClient(conn)
	resp, err := client.CreateItem(ctx, &proto.CreateItemReq{Item: &proto.Item{ItemName: "TestItem1"}})
	require.NoError(t, err)

	assert.Equal(t, resp.Item.ItemName, "TestItem1")
}