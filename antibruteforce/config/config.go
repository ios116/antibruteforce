package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Settings struct {
	LoginRequests    int   `env:"LOGIN_REQUESTS" envDefault:"10"`
	PasswordRequests int   `env:"PASSWORD_REQUESTS" envDefault:"100"`
	IpRequests       int   `env:"IP_REQUESTS" envDefault:"1000"`
	Duration         int32 `env:"DURATION" envDefault:"60"`
}

func NewSettings() *Settings {
	c := &Settings{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}
