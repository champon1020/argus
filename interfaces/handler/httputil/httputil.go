package httputil

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func ParsePage(c echo.Context) (int, error) {
	p := c.QueryParam("p")
	page, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(page), nil
}
