package grpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"integration_test/config"
	"testing"
)

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": t.Token,
	}, nil
}

func (t *tokenAuth) RequireTransportSecurity() bool {
	return false
}

func TestGrpc(t *testing.T) {
	option := grpc.WithPerRPCCredentials(&tokenAuth{"Bearer secret"})
	conf := config.NewGrpcConf()
	address := fmt.Sprintf("%s:%d", conf.GrpcHost, conf.GrpcPort)
	conn, err := grpc.Dial(address, option, grpc.WithInsecure())
	t.Log(address)
	if err != nil {
		t.Fatal("Can't connect to GRPC: ", address)
	}

	server := NewAntiBruteForceClient(conn)
	ctx := context.Background()

	t.Run("Add IP", func(t *testing.T) {
		req := &AddIpRequest{
			Net:  "127.0.0.1/32",
			List: List_BLACK,
		}
		status, err := server.AddIP(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if !status.Ok {
			t.Fatal("status should by true")
		}
	})

	t.Run("Get subnet", func(t *testing.T) {
		req := &GetSubnetRequest{
			Net: "127.0.0.1/32",
		}
		resp, err := server.GetSubnet(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Nets) == 0 {
			t.Fatal("nets should be 1")
		}
		t.Log(resp.Nets[0].Net, resp.Nets[0].List)
	})

	t.Run("Delete ip", func(t *testing.T) {
		req := &DeleteIpRequest{
			Net: "127.0.0.1/32",
		}
		status, err := server.DeleteIP(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if !status.Ok {
			t.Fatal("status should by true")
		}
	})

	t.Run("Check request", func(t *testing.T) {
		req := &CheckRequest{
			Login:    "Admin",
			Password: "123456",
			Ip:       "91.245.34.189",
		}
		for i:=0; i<1100; i++ {
			status, err := server.Check(ctx, req)
			t.Log("===> ",i,status, err)
			//if err != nil {
			//	t.Fatal(err)
			//}
			//if !status.Ok {
			//	t.Fatal("status should be true but get false")
			//}
		}
	})
}
