package httphandler

import (
	"errors"
	"net/http"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/service"
	"github.com/go-playground/validator/v10"

	//"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type authHandler struct {
	//config      *config.AppConfig
	//loc         *time.Location
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler { // config *config.AppConfig, loc *time.Location,
	return &authHandler{
		authService: authService,
	}
}

// POST /api/v1/register
func (h *authHandler) Register(c echo.Context) error {

	ctx := c.Request().Context()

	req, err := buildRegisterRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.AppError{
			Code:    model.StatusBadRequest,
			Message: err.Error(),
		})
	}

	response, err := h.authService.Register(ctx, req)
	if err != nil {
		if err.Error() == model.EmailAlreadyExistMessage {
			return c.JSON(http.StatusOK, model.AppError{
				Code:    model.StatusEmailAlreadyExist,
				Message: "Email already exist, please use different email.",
			})
		}
		if err.Error() == model.ServiceUnavailableMessage {
			return c.JSON(http.StatusOK, model.AppError{
				Code:    model.StatusServiceUnavailable,
				Message: "Service is available between 06:00 and 23:00, please use the service during the operate time.",
			})
		}
		return c.JSON(http.StatusOK, model.AppError{
			Code: model.StatusGenericError,
			Message: "Generic error",
		})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *authHandler) Login(c echo.Context) error {
	return nil
}

func buildRegisterRequest(c echo.Context) (*model.RegisterRequest, error) {

	req := new(model.RegisterRequest) //var req model.RegisterRequest

	if err := c.Bind(req); err != nil { // &req
		return nil, errors.New("Invalid request body")
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func buildLoginRequest(c echo.Context) (*model.LoginRequest, error) {
	//var req model.RegisterRequest
	req := new(model.LoginRequest)

	if err := c.Bind(req); err != nil { // &req
		return nil, errors.New("Invalid request body")
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
