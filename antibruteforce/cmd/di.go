package cmd

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/store/bucketstore"
	"antibruteforce/internal/store/ipstore"
	"antibruteforce/internal/usecase/bucketusecase"

	"go.uber.org/dig"
)

//var ddd func (store entities.BucketStoreManager, settings *config.Settings) *bucketusecase.BucketService {
//
//}

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(config.NewSettings)
	container.Provide(config.NewGrpcConf)
	container.Provide(config.NewDateBaseConf)
	container.Provide(config.CreateLogger)
	container.Provide(config.DBConnection)
	container.Provide(ipstore.NewDbRepo)
	container.Provide(bucketstore.NewBucketStore)
	container.Provide(bucketusecase.NewBucketService)
	return container
}
