package grpcserver

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/usecase"
	"antibruteforce/internal/usecase/bucketusecase"
	"antibruteforce/internal/usecase/ipusecase"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

// RPCServer grpc api
type RPCServer struct {
	Conf           *config.GrpcConf
	Logger         *zap.Logger
	IPService      ipusecase.IPUseCase
	BucketService  bucketusecase.BucketsUseCase
	IntegratorService usecase.Interactor
}

// Check grpc method for check request
func (r *RPCServer) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*StatusResponse, error) {

	req := &entities.Request{
		IP:       in.Ip,
		Login:    in.Login,
		Password: in.Password,
	}
	status, err := r.IntegratorService.CheckRequest(req)
	if err != nil || !status {
		return &StatusResponse{Status: false}, err
	}
	return &StatusResponse{Status: true}, nil
}

// ResetBucket grpc method for reset bucket
func (r *RPCServer) ResetBucket(ctx context.Context, in *ResetBucketRequest, opts ...grpc.CallOption) (*StatusResponse, error) {

	panic("not implementation")
}

// AddIP grpc method for add ip to whitelist or blacklist
func (r *RPCServer) AddIP(ctx context.Context, in *AddIpRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	_, netIP, err := net.ParseCIDR(in.Net)
	if err != nil {
		return &StatusResponse{Status: false}, err
	}
	var kind entities.IPKind
	switch in.List {
	case List_BLACK:
		kind = entities.Black
	case List_WHITE:
		kind = entities.White
	}
	ip := &entities.IPItem{
		Kind: kind,
		IP:   netIP,
	}
	if err := r.IPService.AddIpToList(ctx, ip); err != nil {
		return &StatusResponse{Status: false}, err
	}
	return &StatusResponse{Status: true}, nil
}

// DeleteIP delete ip from list
func (r *RPCServer) DeleteIP(ctx context.Context, in *DeleteIpRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	_, netIP, err := net.ParseCIDR(in.Net)
	if err != nil {
		return &StatusResponse{Status: false}, err
	}
	err = r.IPService.DeleteByIP(ctx, netIP)
	if err != nil {
		return &StatusResponse{Status: false}, err
	}
	return &StatusResponse{Status: true}, nil
}
