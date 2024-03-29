package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

// GrpcConf struct for grpc env
type GrpcConf struct {
	GrpcPort  int    `env:"GRPC_PORT" envDefault:"50051"`
	GrpcHost  string `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	GrpcToken string `env:"GRPC_TOKEN" envDefault:"secret"`
}

// NewGrpcConf parsing grpc env
func NewGrpcConf() *GrpcConf {
	c := &GrpcConf{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}
