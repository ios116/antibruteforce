package cmd

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/store/ipstore"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(config.NewSettings())
	container.Provide(config.NewGrpcConf())
	container.Provide(config.NewDateBaseConf())
	container.Provide(ipstore.NewDbRepo)
	return container
}
