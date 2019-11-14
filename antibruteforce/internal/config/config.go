package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

// Settings are used to set restrictions on requests.
type Settings struct {
	LoginRequests    int   `env:"LOGIN_REQUESTS" envDefault:"10"`
	PasswordRequests int   `env:"PASSWORD_REQUESTS" envDefault:"100"`
	IPRequests       int   `env:"IP_REQUESTS" envDefault:"1000"`
	Duration         int32 `env:"DURATION" envDefault:"60"`
}

// NewSettings new instance for Settings
func NewSettings() *Settings {
	c := &Settings{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}
