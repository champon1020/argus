package handler

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/champon1020/argus/v2"
	"github.com/gin-gonic/gin"
)

var (
	errPageIsNotNumber = errors.New("handler.util: Failed to atoi because page is not number")
	errNumIsNotNumber  = errors.New("handler.util: Failed to atoi because num is not number")
	errParamNotFound   = errors.New("handler.util: Query parameter is not found")
)

// ParsePage parses query parameter to get page number.
func ParsePage(ctx *gin.Context, outc chan<- int, errc chan<- error) {
	defer close(outc)
	pStr, ok := ctx.GetQuery("p")
	if !ok {
		// If query parameter 'p' is not found, default page number is 1.
		outc <- 1
		return
	}

	p, err := strconv.Atoi(pStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		err := argus.NewError(errPageIsNotNumber, err).
			AppendValue("page", pStr)
		errc <- err
		return
	}

	outc <- p
}

// ParseNum parses query parameter to get num.
func ParseNum(ctx *gin.Context, outc chan<- int, errc chan<- error) {
	defer close(outc)
	numStr, ok := ctx.GetQuery("num")
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		err := argus.NewError(errParamNotFound, nil)
		errc <- err
		return
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		err := argus.NewError(errNumIsNotNumber, err).
			AppendValue("num", numStr)
		errc <- err
		return
	}

	outc <- num
}

// ParseTitle parses query parameter to get title string.
func ParseTitle(ctx *gin.Context, outc chan<- string, errc chan<- error) {
	defer close(outc)
	title, ok := ctx.GetQuery("title")
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		err := argus.NewError(errParamNotFound, nil).
			AppendValue("param", "title")
		errc <- err
		return
	}

	outc <- title
}

func printMemory(mem runtime.MemStats) {
	runtime.ReadMemStats(&mem)
	fmt.Println(mem.Alloc, mem.TotalAlloc, mem.HeapAlloc, mem.HeapSys)
}
