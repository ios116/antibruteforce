package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Settings struct {
	LoginN    int   `env:"LoginN" envDefault:"10"`
	PasswordM int   `env:"PASSWORD_M" envDefault:"100"`
	IpK       int   `env:"IP_M" envDefault:"1000"`
	Duration  int32 `env:"DURATION" envDefault:"60"`
}

func NewSettings() *Settings {
	c := &Settings{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}
