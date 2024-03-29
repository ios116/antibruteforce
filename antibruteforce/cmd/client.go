package cmd

import (
	"antibruteforce/internal/config"
	"context"
	"fmt"

	"google.golang.org/grpc"
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

func newGrpcConnection(conf *config.GrpcConf) (*grpc.ClientConn, error) {
	token := fmt.Sprintf("Bearer %s", conf.GrpcToken)
	option := grpc.WithPerRPCCredentials(&tokenAuth{token})
	address := fmt.Sprintf("%s:%d", conf.GrpcHost, conf.GrpcPort)
	conn, err := grpc.Dial(address, option, grpc.WithInsecure())
	return conn, err
}
