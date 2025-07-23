package storage

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewStorage(cfg *StorageConfig) (*gorm.DB, error) {
	switch cfg.Driver {
	case "mysql":
		return NewMysqlStorage(cfg.Mysql)
	case "sqlite":
		return NewSqliteStorage(cfg.Sqlite)
	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", cfg.Driver)
	}
}

func NewMysqlStorage(cfg *MysqlConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewSqliteStorage(cfg *SqliteConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
