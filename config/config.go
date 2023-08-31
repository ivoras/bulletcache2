package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	MemcachedPort uint16 `env:"MEMCACHED_PORT" envDefault:"11211"`
}

func InitConfig() (cfg Config, err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}
	err = env.Parse(&cfg)
	return
}
