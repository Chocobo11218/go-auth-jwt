package httphandler

import (
	"errors"
	"net/http"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/service"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type authHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
		validate: validator.New(),
	}
}

// POST /api/v1/register
func (h *authHandler) Register(c echo.Context) error {

	ctx := c.Request().Context()

	logger.Info(ctx, "AuthHandler - Register called")

	req, err := h.buildRegisterRequest(c)
	if err != nil {
		logger.Error(
			ctx, 
			"Register Handler: buildRegisterRequest error", 
			zap.Error(err),
		)
		return c.NoContent(http.StatusBadRequest)
	}
	
	response, err := h.authService.Register(ctx, req)
	if err != nil {

		if err.Error() == model.EmailAlreadyExistMessage {
			logger.Error(ctx, "Register Handler: Email already exist", zap.Error(err))
			return c.JSON(http.StatusOK, model.AppResponse{
				Code:    model.StatusEmailAlreadyExist,
				Message: model.EmailAlreadyExistMessage,
			})
		}
		if err.Error() == model.ServiceUnavailableMessage {
			logger.Error(ctx, "Register Handler: Service hour restricted", zap.Error(err))
			return c.JSON(http.StatusOK, model.AppResponse{
				Code:    model.StatusServiceUnavailable,
				Message: model.ServiceUnavailableMessage,
			})
		}
		logger.Error(ctx, "Register Handler: Generic error", zap.Error(err))
		return c.JSON(http.StatusOK, model.AppResponse{
			Code: model.StatusGenericError,
			Message: model.GenericErrorMessage,
		})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *authHandler) Login(c echo.Context) error {
	
	ctx := c.Request().Context()

	logger.Info(ctx, "AuthHandler - Login called")

	req, err := h.buildLoginRequest(c) 
	if err != nil {
		logger.Error(ctx, "Login Handler: buildLoginRequest error", zap.Error(err))
		return c.NoContent(http.StatusBadRequest)
	}

	response, err := h.authService.Login(ctx, req)
	if err != nil {
		if err.Error() == model.InvalidCredentialMessage {
			logger.Error(ctx, "Login Handler: Invalid credential", zap.Error(err))
			return c.NoContent(http.StatusUnauthorized)
		}
		if err.Error() == model.ServiceUnavailableMessage {
			logger.Error(ctx, "Login Handler: Service hours restricted", zap.Error(err))
			return c.JSON(http.StatusOK, model.AppResponse{
				Code: model.StatusServiceUnavailable,
				Message: model.ServiceUnavailableMessage,
			})
		}
		logger.Error(ctx, "Login Handler: Generic error", zap.Error(err))
		return c.JSON(http.StatusOK, model.AppResponse{
			Code: model.StatusGenericError,
			Message: model.GenericErrorMessage,
		})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *authHandler) buildRegisterRequest(c echo.Context) (*model.RegisterRequest, error) {

	req := new(model.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return nil, errors.New("Invalid request body")
	}
	err := h.validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (h *authHandler) buildLoginRequest(c echo.Context) (*model.LoginRequest, error) {
	//var req model.RegisterRequest
	req := new(model.LoginRequest)

	if err := c.Bind(req); err != nil { // &req
		return nil, errors.New("Invalid request body") 
	}

	err := h.validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
