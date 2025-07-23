//go:build wireinject
// +build wireinject

package user

import (
	"github.com/3c4j/monica/lib/storage"
	"github.com/3c4j/monica/monica"
	"github.com/3c4j/monica/pkg/logger"
	"github.com/3c4j/monica/user/config"
	"github.com/3c4j/monica/user/entity"
	"github.com/3c4j/monica/user/handler"
	"github.com/3c4j/monica/user/service"
	"github.com/google/wire"
)

func InitModule(lg *logger.Logger) (monica.Module, error) {
	panic(wire.Build(
		service.NewAuthService,
		service.NewJwtService,
		handler.NewAuthHandler,
		entity.NewUserRepository,
		config.NewConfig,
		config.NewStorageConfig,
		storage.NewStorage,
		NewModule,
		wire.Bind(new(monica.Module), new(*Module)),
	))
}
