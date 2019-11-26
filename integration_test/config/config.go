package config

import (
	"github.com/caarlos0/env/v6"
	"log"
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
