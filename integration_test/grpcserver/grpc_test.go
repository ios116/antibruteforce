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
	defer conn.Close()
	settings := config.NewSettings()
	t.Log(address)
	if err != nil {
		t.Fatal("Can't connect to GRPC: ", address)
	}

	server := NewAntiBruteForceClient(conn)
	ctx := context.Background()

	t.Run("Check Login", func(t *testing.T) {
		req := &CheckRequest{
			Login:    "Admin",
			Password: "123456",
			Ip:       "91.245.34.189",
		}
		t.Logf("%d request for login per %ds", settings.LoginLimit+1, settings.Duration)
		for i := 0; i < settings.LoginLimit+1; i++ {
			status, err := server.Check(ctx, req)
			if i == settings.LoginLimit && err == nil && status != nil {
				t.Fatalf("Login - %s should be rejected ", req.Login)
			}
		}
	})

	t.Run("Added ip 91.245.34.189 to white list and test by login", func(t *testing.T) {
		ipToWhiteList := &AddIpRequest{
			Net:  "91.245.34.189/32",
			List: List_WHITE,
		}
		_, err := server.AddIP(ctx, ipToWhiteList)
		if err != nil {
			t.Fatal(err)
		}
		req := &CheckRequest{
			Login:    "Admin",
			Password: "123456",
			Ip:       "91.245.34.189",
		}
		t.Logf("%d request for login per %ds", settings.LoginLimit+1, settings.Duration)
		for i := 0; i < settings.LoginLimit+1; i++ {
			status, err := server.Check(ctx, req)
			if  err != nil || status == nil {
				t.Fatalf("Login - %s should be rejected ", req.Login)
			}
		}
	})

	t.Run("Get subnet", func(t *testing.T) {
		req := &GetSubnetRequest{
			Net: "91.245.34.189/32",
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

	t.Run("Delete subnet", func(t *testing.T) {
		req := &DeleteIpRequest{
			Net: "91.245.34.189/32",
		}
		status, err := server.DeleteIP(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if !status.Ok {
			t.Fatal("status should by true")
		}
	})

	t.Run("Check brute by IP", func(t *testing.T) {
		t.Logf("%d request for login per %ds", settings.IPLimit+1, settings.Duration)
		for i := 0; i < settings.LoginLimit+1; i++ {
			login := fmt.Sprintf("Login-%d", i)
			password := fmt.Sprintf("Password-%d", i)
			req := &CheckRequest{
				Login:    login,
				Password: password,
				Ip:       "91.245.34.12",
			}
			status, _ := server.Check(ctx, req)
			if i == settings.IPLimit && err == nil && status != nil {
				t.Fatalf("IP - %s should be rejected ", req.Ip)
			}
		}
	})
}
