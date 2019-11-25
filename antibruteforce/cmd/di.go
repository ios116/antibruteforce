package cmd

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/grpcserver"
	"antibruteforce/internal/store/bucketstore"
	"antibruteforce/internal/store/ipstore"
	"antibruteforce/internal/usecase/bucketusecase"
	"antibruteforce/internal/usecase/interactor"
	"antibruteforce/internal/usecase/ipusecase"
	"log"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func bucketService(store *bucketstore.BucketStore, settings *config.Settings, logger *zap.Logger) *bucketusecase.BucketService {
	return bucketusecase.NewBucketService(store, settings, logger)
}

func ipService(settings *config.Settings, IPStore *ipstore.DbRepo) *ipusecase.IPService {
	return ipusecase.NewIPService(settings, IPStore)
}

func ineractor(ipService *ipusecase.IPService, bucketService *bucketusecase.BucketService) *interactor.Connector {
	return interactor.NewConnector(ipService, bucketService)
}

func grpcCast(conf *config.GrpcConf, logger *zap.Logger, IPService *ipusecase.IPService, bucketService *bucketusecase.BucketService, integratorService *interactor.Connector) *grpcserver.RPCServer {
	return grpcserver.NewRPCServer(conf, logger, IPService, bucketService, integratorService)
}

// BuildContainer creates dependency injection
func BuildContainer() *dig.Container {
	container := dig.New()

	// create config
	if err := container.Provide(config.NewSettings); err != nil {
		log.Println(err)
	}
	if err := container.Provide(config.NewGrpcConf); err != nil {
		log.Println(err)
	}
	if err := container.Provide(config.NewDateBaseConf); err != nil {
		log.Println(err)
	}
	if err := container.Provide(config.CreateLogger); err != nil {
		log.Println(err)
	}
	if err := container.Provide(config.DBConnection); err != nil {
		log.Println(err)
	}
	// create ip store
	if err := container.Provide(ipstore.NewDbRepo); err != nil {
		log.Println(err)
	}

	// create bucket store
	if err := container.Provide(bucketstore.NewBucketStore); err != nil {
		log.Println(err)
	}

	// create new bucket service
	if err := container.Provide(bucketService); err != nil {
		log.Println(err)
	}

	// create new ip service
	if err := container.Provide(ipService); err != nil {
		log.Println(err)
	}

	// create integrator service
	if err := container.Provide(ineractor); err != nil {
		log.Println(err)
	}

	// create grpc server
	if err := container.Provide(grpcCast); err != nil {
		log.Println(err)
	}

	return container
}
