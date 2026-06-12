package httphandler

import (
	"errors"
	"net/http"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/middleware"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/model"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/service"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler interface {
	GetMe(c echo.Context) error
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

// GET /api/v1/me
func (h *userHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()

	logger.Info(ctx, "UserHandler - GetMe called")

	userID, ok := c.Get(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		return c.NoContent(http.StatusUnauthorized)
	}

	response, err := h.userService.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, errors.New(model.UserNotFoundMessage)) || err.Error() == model.UserNotFoundMessage {
			logger.Error(ctx, "UserHandler - GetMe: user not found", zap.Error(err))
			return c.NoContent(http.StatusNotFound)
		}
		logger.Error(ctx, "UserHandler - GetMe: generic error", zap.Error(err))
		return c.JSON(http.StatusOK, model.AppResponse{
			Code:    model.StatusGenericError,
			Message: model.GenericErrorMessage,
		})
	}

	return c.JSON(http.StatusOK, response)
}
