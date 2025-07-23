package entity

import (
	"github.com/3c4j/monica/lib/repository"
	"github.com/3c4j/monica/pkg/logger"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
}

func NewOfflineUser(id uint, username string) *User {
	return &User{
		Model: gorm.Model{
			ID: id,
		},
		Username: username,
	}
}

func NewUserRepository(db *gorm.DB, lg *logger.Logger) *repository.Repository[User] {
	return repository.NewRepository[User](db)
}
