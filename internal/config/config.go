package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type RootConfig struct {
	Server struct {
		Port int `toml:"port"`
		Jwt  struct {
			PrivateKey string `toml:"privateKey"`
			PublicKey  string `toml:"publicKey"`
		} `toml:"jwt"`
	} `toml:"server"`
	Database struct {
		URI string `toml:"uri"`
	} `toml:"database"`
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
