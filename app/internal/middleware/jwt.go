package middleware

import (
	"net/http"
	"strings"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const UserIDKey = "user_id"

func JWTAuth(conf *config.AppConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.NoContent(http.StatusUnauthorized)
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(conf.Secret.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}

			c.Set(UserIDKey, uint(userIDFloat))
			return next(c)
		}
	}
}
