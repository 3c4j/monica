package config

import (
	"github.com/3c4j/monica/lib/storage"
	"github.com/spf13/viper"
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := viper.UnmarshalKey("user", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewStorageConfig(cfg *Config) *storage.StorageConfig {
	return cfg.Storage
}

type HttpConfig struct {
	Port           int    `mapstructure:"port"`
	Host           string `mapstructure:"host"`
	ReadTimeout    int    `mapstructure:"read_timeout"`
	WriteTimeout   int    `mapstructure:"write_timeout"`
	IdleTimeout    int    `mapstructure:"idle_timeout"`
	MaxHeaderBytes int    `mapstructure:"max_header_bytes"`
}

type JwtConfig struct {
	Secret string `mapstructure:"secret"`
}

type Config struct {
	Storage *storage.StorageConfig `mapstructure:"storage"`
	Http    *HttpConfig            `mapstructure:"http"`
	Jwt     *JwtConfig             `mapstructure:"jwt"`
}
