package handler

import (
	"net/http"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/interfaces/auth"
	"github.com/champon1020/argus/util"
	"github.com/labstack/echo/v4"
)

// AuthHandler is handler interface for authentication.
type AuthHandler interface {
	VerifyToken(c echo.Context) error
}

type authHandler struct {
	logger *argus.Logger
}

// NewAuthHandler creates authHandler.
func NewAuthHandler(logger *argus.Logger) AuthHandler {
	return &authHandler{logger: logger}
}

func (aH *authHandler) VerifyToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token, err := util.ExtractBearerToken(authHeader)
	if err != nil {
		aH.logger.Error(c, http.StatusBadRequest, err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid authorization header")
	}

	if err := auth.VerifyJWTToken(token); err != nil {
		aH.logger.Error(c, http.StatusForbidden, err)
		return echo.NewHTTPError(http.StatusForbidden, "failed to verify token")
	}

	return c.String(http.StatusOK, "success")
}
