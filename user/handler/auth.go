package handler

import (
	"net/http"

	"github.com/3c4j/monica/pkg/logger"
	"github.com/3c4j/monica/user/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtService  *service.JwtService
	logger      *logger.Logger
}

func NewAuthHandler(authService *service.AuthService, jwtService *service.JwtService, lg *logger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, jwtService: jwtService, logger: lg.With(logger.F{"module": "user.handler.auth"})}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("login: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		h.logger.Errorf("login: %s", err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, err := h.jwtService.GenerateToken(ctx, user, "login")
	if err != nil {
		h.logger.Errorf("login: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Confirm  string `json:"confirm"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("register: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Password != req.Confirm {
		h.logger.Errorf("register: password and confirm password do not match")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password and confirm password do not match"})
		return
	}
	user, err := h.authService.Register(ctx, req.Username, req.Password)
	if err != nil {
		h.logger.Errorf("register: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := h.jwtService.GenerateToken(ctx, user, "register")
	if err != nil {
		h.logger.Errorf("register: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
