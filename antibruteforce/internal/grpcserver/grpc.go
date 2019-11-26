package grpcserver

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/domain/entities"
	"antibruteforce/internal/usecase/bucketusecase"
	"antibruteforce/internal/usecase/interactor"
	"antibruteforce/internal/usecase/ipusecase"
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RPCServer grpc api
type RPCServer struct {
	Conf              *config.GrpcConf
	Logger            *zap.Logger
	IPService         ipusecase.IPUseCase
	BucketService     bucketusecase.BucketsUseCase
	IntegratorService interactor.ConnectorUseCase
}

// NewRPCServer constructor for GRPC server
func NewRPCServer(conf *config.GrpcConf, logger *zap.Logger, IPService ipusecase.IPUseCase, bucketService bucketusecase.BucketsUseCase, integratorService interactor.ConnectorUseCase) *RPCServer {
	return &RPCServer{Conf: conf, Logger: logger, IPService: IPService, BucketService: bucketService, IntegratorService: integratorService}
}

// Start - init RPC server
func (r *RPCServer) Start() {
	address := fmt.Sprintf("%s:%d", r.Conf.GrpcHost, r.Conf.GrpcPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		r.Logger.Fatal("Cannot start RPC server", zap.String("err", err.Error()))
	}
	// server := grpc.NewServer(grpc.UnaryInterceptor(newInterceptor(g.logger, g.conf.GrpcToken)))
	go r.BucketService.BucketCollector()
	server := grpc.NewServer()
	RegisterAntiBruteForceServer(server, r)
	r.Logger.Info("Starting RPC server", zap.String("address", address))
	err = server.Serve(lis)
	if err != nil {
		r.Logger.Fatal("Cannot start listen port", zap.String("err", err.Error()))
	}
}

// Check grpc method for check request
func (r *RPCServer) Check(ctx context.Context, in *CheckRequest) (*StatusResponse, error) {
	req := &entities.Request{
		IP:       in.Ip,
		Login:    in.Login,
		Password: in.Password,
	}
	status, err := r.IntegratorService.CheckRequest(req)
	return &StatusResponse{Ok: status}, err
}

// ResetBucket grpc method for reset bucket
func (r *RPCServer) ResetBucket(ctx context.Context, in *ResetBucketRequest) (*StatusResponse, error) {
	var kind entities.KindBucket
	switch in.Kind {
	case BucketKind_IP:
		kind = entities.IP
	case BucketKind_LOGIN:
		kind = entities.Login
	case BucketKind_PASSWORD:
		kind = entities.Password
	}
	hash := entities.Hash{
		Kind: kind,
		Key:  in.Key,
	}
	err := r.BucketService.ResetBucket(hash)
	return &StatusResponse{Ok: err == nil}, err
}

// AddIP grpc method for add ip to whitelist or blacklist
func (r *RPCServer) AddIP(ctx context.Context, in *AddIpRequest) (*StatusResponse, error) {
	_, netIP, err := net.ParseCIDR(in.Net)
	if err != nil {
		return &StatusResponse{Ok: false}, err
	}
	var kind entities.IPKind
	switch in.List {
	case List_BLACK:
		kind = entities.Black
	case List_WHITE:
		kind = entities.White
	}
	ip := &entities.IPListRow{
		Kind: kind,
		IP:   netIP,
	}
	err = r.IPService.AddNet(ctx, ip)
	return &StatusResponse{Ok: err == nil}, err
}

// DeleteIP delete ip from list
func (r *RPCServer) DeleteIP(ctx context.Context, in *DeleteIpRequest) (*StatusResponse, error) {
	_, netIP, err := net.ParseCIDR(in.Net)
	if err != nil {
		return &StatusResponse{Ok: false}, err
	}
	err = r.IPService.DeleteNet(ctx, netIP)
	return &StatusResponse{Ok: err == nil}, err
}

// GetSubnet get all subnet by net
func (r *RPCServer) GetSubnet(ctx context.Context, in *GetSubnetRequest) (*GetSubnetResponse, error) {
	results, err := r.IPService.GetSubnet(ctx, in.Net)
	if err != nil {
		return nil, err
	}
	var nets []*Net
	for _, item := range results {
		var list List
		switch item.Kind {
		case entities.Black:
			list = List_BLACK
		case entities.White:
			list = List_WHITE
		}
		nets = append(nets, &Net{
			Net:  item.IP.String(),
			List: list,
		})
	}
	return &GetSubnetResponse{
		Nets: nets,
	}, nil
}
