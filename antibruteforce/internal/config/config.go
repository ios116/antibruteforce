package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// Settings are used to set restrictions on LIMIT.
type Settings struct {
	LoginLimit    int    `env:"LOGIN_LIMIT" envDefault:"10"`
	PasswordLimit int    `env:"PASSWORD_LIMIT" envDefault:"100"`
	IPLimit       int    `env:"IP_LIMIT" envDefault:"1000"`
	Duration      int32  `env:"DURATION" envDefault:"60"`
	Build         string `env:"BUILD" envDefault:"dev"`
}

// NewSettings new instance for Settings
func NewSettings() *Settings {
	c := &Settings{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}

// CreateLogger creating the logger
func CreateLogger(c *Settings) (logger *zap.Logger, err error) {
	if c.Build == "dev" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	return logger, err
}
