package httphandler

import (
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type authHandler struct {
	config      *config.AppConfig
	loc         *time.Location
	authService service.AuthService
}

func NewAuthHandler(
	config *config.AppConfig, 
	loc *time.Location, 
	authService service.AuthService,
) AuthHandler{
	return &authHandler{
		config: config,
		loc: loc,
		authService: authService,
	}
}

func (h *authHandler) Register(c echo.Context) error {
	return nil
}

func (h *authHandler) Login(c echo.Context) error {
	return nil
}