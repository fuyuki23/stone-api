package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type RootConfig struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
}

var instance RootConfig

func Load() error {
	if err := cleanenv.ReadConfig("config.toml", &instance); err != nil {
		return err
	}

	return nil
}

func Get() RootConfig {
	return instance
}
