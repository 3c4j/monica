package storage

import "fmt"

type StorageConfig struct {
	Driver string        `mapstructure:"driver"`
	Mysql  *MysqlConfig  `mapstructure:"mysql"`
	Sqlite *SqliteConfig `mapstructure:"sqlite"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

func (c *MysqlConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DB)
}

type SqliteConfig struct {
	Path string `mapstructure:"path"`
}
