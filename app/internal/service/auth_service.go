package service

import (
	"context"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	//"github.com/labstack/echo/v4"
)

type AuthService interface {
	Register(ctx context.Context) error
	Login(ctx context.Context) error
}

type authService struct {
	config   *config.AppConfig
	userRepo repository.UserRepository
	// internalKafka         repository.InternalKafka
	loc *time.Location
}

func NewAuthService(
	config *config.AppConfig,
	userRepo repository.UserRepository,
	// internalKafka         repository.InternalKafka,
	loc *time.Location,
) AuthService {
	return &authService{
		config:   config,
		userRepo: userRepo,
		// internalKafka:  internalKafka,
		loc: loc,
	}
}

func (s *authService) Register(ctx context.Context) error {
	return nil
}

func (s *authService) Login(ctx context.Context) error {
	return nil
}
