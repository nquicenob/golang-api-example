package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Debug    bool   `envconfig:"DEBUG"`
	Port     string `envconfig:"PORT"`
	LogLevel string `envconfig:"LOG_LEVEL"`
	RetryTx  int    `envconfig:"MAX_RETRY_TX"`

	IsHeroku bool `envconfig:"IS_HEROKU"`

	DBHost           string `envconfig:"DATABASE_HOST"`
	DBPort           string `envconfig:"DATABASE_PORT"`
	DBSSLMode        string `envconfig:"DATABASE_SSL_MODE"`
	DBUser           string `envconfig:"POSTGRES_USER"`
	DBName           string `envconfig:"POSTGRES_DB"`
	DBPass           string `envconfig:"POSTGRES_PASSWORD"`
	DBConnRetries    int    `envconfig:"DATABASE_CONN_MAX_RETRIES"`
	DBConnectTimeout string `envconfig:"DATABASE_CONN_TIMEOUT"`
}

func New() (s *Specification, err error) {
	s = new(Specification)
	if err := envconfig.Process("golangapi", s); err != nil {
		return s, err
	}
	return s, nil
}
