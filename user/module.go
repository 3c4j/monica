package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/3c4j/monica/pkg/logger"
	"github.com/3c4j/monica/user/config"
	"github.com/3c4j/monica/user/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Module struct {
	Config      *config.Config
	AuthHandler *handler.AuthHandler
	server      *http.Server
	logger      *logger.Logger
}

func NewModule(lg *logger.Logger, cfg *config.Config, authHandler *handler.AuthHandler) (*Module, error) {
	return &Module{Config: cfg, AuthHandler: authHandler, logger: lg.With(logger.F{"module": "user"})}, nil
}

func (m *Module) Run(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.POST("/auth/login", m.AuthHandler.Login)
	router.POST("/auth/register", m.AuthHandler.Register)
	m.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", m.Config.Http.Host, m.Config.Http.Port),
		Handler: router,
	}
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.logger.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	return m.Shutdown(ctx)
}

func (m *Module) Shutdown(ctx context.Context) error {
	return m.server.Close()
}
