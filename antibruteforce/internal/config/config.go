package config

import (
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"log"
)

// Settings are used to set restrictions on requests.
type Settings struct {
	LoginRequests    int    `env:"LOGIN_REQUESTS" envDefault:"10"`
	PasswordRequests int    `env:"PASSWORD_REQUESTS" envDefault:"100"`
	IPRequests       int    `env:"IP_REQUESTS" envDefault:"1000"`
	Duration         int32  `env:"DURATION" envDefault:"60"`
	Build            string `env:"BUILD" envDefault:"dev"`
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
