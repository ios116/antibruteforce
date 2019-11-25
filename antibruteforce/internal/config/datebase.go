package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"

	// postgres config
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// DateBaseConf config
type DateBaseConf struct {
	BdPassword string `env:"POSTGRES_PASSWORD" envDefault:"123456"`
	BdUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	BdHost     string `env:"POSTGRES_HOST" envDefault:"0.0.0.0"`
	BdName     string `env:"POSTGRES_DB" envDefault:"force"`
}

// NewDateBaseConf parsing env
func NewDateBaseConf() *DateBaseConf {
	c := &DateBaseConf{}
	if err := env.Parse(c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return c
}

// DBConnection - connection for BD
// postgres://myuser:mypass@localhost:5432/mydb?sslmode=verify­full
// export POSTGRESQL_URL=postgres://calendar:123456@localhost:5432/calendar?sslmode=disable
func DBConnection(c *DateBaseConf) (*sqlx.DB, error) {
	var params = fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable", c.BdUser, c.BdName, c.BdHost, c.BdPassword)
	db, err := sqlx.Connect("pgx", params)
	return db, err
}
