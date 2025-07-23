package service

import (
	"context"
	"errors"
	"time"

	"github.com/3c4j/monica/pkg/logger"
	"github.com/3c4j/monica/user/config"
	"github.com/3c4j/monica/user/entity"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtService struct {
	secret string
	logger *logger.Logger
}

func NewJwtService(cfg *config.Config, lg *logger.Logger) *JwtService {
	return &JwtService{secret: cfg.Jwt.Secret, logger: lg.With(logger.F{"module": "user.service.jwt"})}
}

func (svc *JwtService) GenerateToken(ctx context.Context, user *entity.User, loginType string) (string, error) {
	expiresIn := time.Hour * 24
	if loginType == "login" {
		expiresIn = time.Hour * 48
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"exp":      time.Now().Add(expiresIn).Unix(),
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
		"iss":      "monica." + loginType,
		"aud":      "monica.auth",
		"jti":      uuid.New().String(),
		"username": user.Username,
	})
	return token.SignedString([]byte(svc.secret))
}

func (svc *JwtService) ParseToken(ctx context.Context, token string) (jwt.MapClaims, error) {
	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(svc.secret), nil
	})
	if err != nil {
		svc.logger.Errorf("parse token: %s", err.Error())
		return nil, err
	}
	claimsMap, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		svc.logger.Errorf("parse token: invalid token")
		return nil, errors.New("invalid token")
	}
	return claimsMap, nil
}

func (svc *JwtService) ValidateToken(ctx context.Context, token string) bool {
	claims, err := svc.ParseToken(ctx, token)
	if err != nil {
		svc.logger.Errorf("validate token: %s", err.Error())
		return false
	}
	return claims.Valid() == nil
}
