package service

import (
	"context"
	"errors"

	"github.com/3c4j/monica/lib/repository"
	"github.com/3c4j/monica/pkg/logger"
	"github.com/3c4j/monica/pkg/utils"
	"github.com/3c4j/monica/user/entity"
	"gorm.io/gorm"
)

type AuthService struct {
	repo   *repository.Repository[entity.User]
	logger *logger.Logger
}

func NewAuthService(repo *repository.Repository[entity.User], lg *logger.Logger) *AuthService {
	return &AuthService{repo: repo, logger: lg.With(logger.F{"module": "user.service.auth"})}
}

func (svc *AuthService) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := svc.repo.FindOne(ctx, entity.User{Username: username})
	if err != nil {
		svc.logger.Errorf("login: %s", err.Error())
		return nil, err
	}
	if !utils.ComparePassword(user.Password, password) {
		svc.logger.Errorf("login: invalid password")
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func (svc *AuthService) Register(ctx context.Context, username, password string) (*entity.User, error) {
	exists, err := svc.repo.FindOne(ctx, entity.User{Username: username})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		svc.logger.Errorf("register: %s", err.Error())
		return nil, err
	}
	if exists != nil {
		svc.logger.Errorf("register: username already exists")
		return nil, errors.New("username already exists")
	}
	user := &entity.User{
		Username: username,
		Password: utils.HashPassword(password),
	}
	err = svc.repo.Create(ctx, user)
	if err != nil {
		svc.logger.Errorf("register: %s", err.Error())
		return nil, err
	}
	return user, nil
}
