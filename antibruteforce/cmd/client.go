package cmd

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/grpcserver"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
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

var grpcClient = &cobra.Command{
	Use:   "client",
	Short: "client grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		container := BuildContainer()
		err := container.Invoke(func(serverGRPS *grpcserver.RPCServer) {
			option := grpc.WithPerRPCCredentials(&tokenAuth{"Bearer secret"})
			conf := config.NewGrpcConf()
			address := fmt.Sprintf("%s:%d", conf.GrpcHost, conf.GrpcPort)
			conn, err := grpc.Dial(address, option, grpc.WithInsecure())
			log.Println(address)
			if err != nil {
				log.Fatal(err)
			}
			server := grpcserver.NewAntiBruteForceClient(conn)
			ctx := context.Background()
			req := &grpcserver.CheckRequest{
				Login:    "Admin",
				Password: "123456",
				Ip:       "91.245.34.189",
			}
			for i:=0; i<10; i++ {
				status, err := server.Check(ctx, req)
				log.Println("===> ",i,status, err)
			}
		})
		if err != nil {
			log.Println(err)
		}
	},
}
