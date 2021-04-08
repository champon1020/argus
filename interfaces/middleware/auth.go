package middleware

import (
	"net/http"

	"github.com/champon1020/argus/interfaces/auth"
	"github.com/champon1020/argus/util"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware verifies jwt token.
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		token, err := util.ExtractBearerToken(authHeader)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := auth.VerifyJWTToken(token); err != nil {
			return c.String(http.StatusForbidden, err.Error())
		}

		return next(c)
	}
}
