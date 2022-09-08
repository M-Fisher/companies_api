package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                string        `envconfig:"PORT" default:":8080"`
	IPApiRequestTimeout time.Duration `envconfig:"IPAPI_REQUEST_TIMEOUT" default:"2s"`
	Postgres            DB            `envconfig:"POSTGRES"`
	Kafka               Kafka         `envconfig:"KAFKA"`
	DevMode             bool          `envconfig:"DEVELOPMENT_MODE" default:"false"`
	JWTSecret           string        `envconfig:"JWT_SECRET" default:"test"`
}

type DB struct {
	User            string        `envconfig:"USER" required:"true"`
	Password        string        `envconfig:"PASSWORD" required:"true"`
	Host            string        `envconfig:"HOST" required:"true"`
	Database        string        `envconfig:"DATABASE" required:"true"`
	MaxIdleConnTime time.Duration `envconfig:"MAX_IDLE_CONN_TIME" default:"5m"`
	MaxConns        int           `envconfig:"MAX_CONNS" default:"20"`
	ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"10m"`
}

type Kafka struct {
	Topic        string        `envconfig:"TOPIC" default:"companies_update"`
	Host         string        `envconfig:"HOST" required:"true"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"2s"`
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess("", &c)
	return &c
}
