package handler

import (
	"errors"
	"strconv"

	"github.com/champon1020/argus/v2"
	"github.com/gin-gonic/gin"
)

var (
	errPageIsNotNumber = errors.New("handler.util: Failed to atoi because page is not number")
	errNumIsNotNumber  = errors.New("handler.util: Failed to atoi because num is not number")
)

// ParsePage parses query parameter to get page number.
func ParsePage(c *gin.Context) (p int, err error) {
	// If query parameter 'p' is not found, default page number is 1.
	p = 1

	if pStr, ok := c.GetQuery("p"); ok {
		if p, err = strconv.Atoi(pStr); err != nil {
			err = argus.NewError(errPageIsNotNumber, err).
				AppendValue("page", pStr)
			return
		}
	}

	return
}

// ParseNum parses query parameter to get num.
func ParseNum(c *gin.Context) (num int, err error) {
	if numStr, ok := c.GetQuery("num"); ok {
		if num, err = strconv.Atoi(numStr); err != nil {
			err = argus.NewError(errNumIsNotNumber, err).
				AppendValue("num", numStr)
			return
		}
	}
	return
}
