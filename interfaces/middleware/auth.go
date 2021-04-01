package middleware

import (
	"errors"
	"strings"

	"github.com/champon1020/argus/service"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware verifies jwt token.
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		token, err := extractToken(auth)
		if err != nil {
			return err
		}

		if err := service.VerifyJWTToken(token); err != nil {
			return err
		}

		return next(c)
	}
}

func extractToken(auth string) (string, error) {
	el := strings.Split(auth, "Bearer ")
	if len(el) < 2 {
		return "", errors.New("invalid authorization header")
	}

	return el[1], nil
}
