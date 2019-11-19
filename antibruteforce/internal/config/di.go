package config

import (
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
	container.Provide(NewSettings)
	container.Provide(NewGrpcConf)
	container.Provide(NewDateBaseConf)
	container.Provide(CreateLogger)
	container.Provide(DBConnection)
	container.Provide(ipstore.NewDbRepo)
	container.Provide(bucketstore.NewBucketStore)
    container.Provide( bucketusecase.NewBucketService)
	return container
}