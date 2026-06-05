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

	req, err := buildRegisterRequest(c, h.validate)
	if err != nil {
		logger.Error(
			ctx, 
			"AuthHandler - buildRegisterRequest error", 
			zap.Error(err),
		)
		
		return c.JSON(http.StatusBadRequest, model.AppResponse{
			Code:    model.StatusBadRequest,
			Message: model.GenericErrorMessage,  // modify
		})
	}

	response, err := h.authService.Register(ctx, req)
	if err != nil {

		if err.Error() == model.EmailAlreadyExistMessage {
			logger.Error(ctx, "AuthHandler - Register email already exist", zap.Error(err))
			return c.JSON(http.StatusOK, model.AppResponse{
				Code:    model.StatusEmailAlreadyExist,
				Message: model.EmailAlreadyExistMessage,
			})
		}
		if err.Error() == model.ServiceUnavailableMessage {
			logger.Error(ctx, "AuthHandler - Service hour restricted", zap.Error(err))
			return c.JSON(http.StatusOK, model.AppResponse{
				Code:    model.StatusServiceUnavailable,
				Message: model.ServiceUnavailableMessage,
			})
		}
		return c.JSON(http.StatusOK, model.AppResponse{
			Code: model.StatusGenericError,
			Message: model.GenericErrorMessage,
		})
	}
	//logger.Info(ctx, "AuthHandler - Register Success", zap.String("email", req.Email))
	return c.JSON(http.StatusOK, response)
}

func (h *authHandler) Login(c echo.Context) error {
	
	ctx := c.Request().Context()

	logger.Info(ctx, "AuthHandler - Login called")

	req, err := buildLoginRequest(c, h.validate) 
	if err != nil {
		logger.Error(ctx, "AuthHandler - buildLoginRequest error", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.AppResponse{
			Code: model.StatusBadRequest,
			Message: model.GenericErrorMessage, 
		})
	}

	response, err := h.authService.Login(ctx, req)
	if err != nil {
		if err.Error() == model.InvalidCredentialMessage { // errors.Is(err, model.)
			logger.Error(ctx, "AuthHandler - Login invalid credential", zap.Error(err))
			return c.JSON(http.StatusOK, model.AppResponse{
				Code: model.StatusInvalidCredential,
				Message: model.InvalidCredentialMessage,
			})
		}
		if err.Error() == model.ServiceUnavailableMessage {
			logger.Error(ctx, "AuthHandler - Service hours restricted", zap.Error(err))
			return c.JSON(http.StatusOK, model.AppResponse{
				Code: model.StatusServiceUnavailable,
				Message: model.ServiceUnavailableMessage,
			})
		}
		return c.JSON(http.StatusOK, model.AppResponse{
			Code: model.StatusGenericError,
			Message: model.GenericErrorMessage,
		})
	}
	//logger.Info(ctx, "AuthHandler - Login Success", zap.String("email", req.Email))
	return c.JSON(http.StatusOK, response)
}

func buildRegisterRequest(c echo.Context, validate *validator.Validate) (*model.RegisterRequest, error) {

	req := new(model.RegisterRequest) //var req model.RegisterRequest

	if err := c.Bind(req); err != nil { // &req
		return nil, errors.New("Invalid request body")
	}

	//validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func buildLoginRequest(c echo.Context, validate *validator.Validate) (*model.LoginRequest, error) {
	//var req model.RegisterRequest
	req := new(model.LoginRequest)

	if err := c.Bind(req); err != nil { // &req
		return nil, errors.New("Invalid request body") 
	}

	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
