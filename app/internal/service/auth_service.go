package service

import (
	"context"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	//"github.com/labstack/echo/v4"
)

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (model.AppResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (model.AppResponse, error)
}

type authService struct {
	config   *config.AppConfig
	userRepo repository.UserRepository
	// internalKafka         repository.InternalKafka
	loc *time.Location
}

func NewAuthService(authRepository repository.UserRepository) AuthService {
	return nil // &authService{}
}

func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (model.AppResponse, error) {
	return model.AppResponse{}, nil
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (model.TokenData, error) {
	return model.TokenData{}, nil
}
