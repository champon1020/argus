package handler

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/champon1020/argus"
	"github.com/gin-gonic/gin"
)

var (
	errParamIsNotNumber = errors.New("handler.util: Failed to atoi because parameter is not number")
	errParamNotFound    = errors.New("handler.util: Query parameter is not found")
)

// ParsePage parses query parameter to get title string.
func ParsePage(ctx *gin.Context, outCh chan<- int, errCh chan<- error) {
	parseIntParam(ctx, "p", outCh, errCh)
}

// ParseNum parses query parameter to get title string.
func ParseNum(ctx *gin.Context, outCh chan<- int, errCh chan<- error) {
	parseIntParam(ctx, "num", outCh, errCh)
}

// ParseTitle parses query parameter to get title string.
func ParseTitle(ctx *gin.Context, outCh chan<- string, errCh chan<- error) {
	parseStringParam(ctx, "title", outCh, errCh)
}

// ParseID parses query parameter to get string id.
func ParseID(ctx *gin.Context, outCh chan<- string, errCh chan<- error) {
	parseStringParam(ctx, "id", outCh, errCh)
}

// ParseCategoryID parses query parameter to get category id.
func ParseCategoryID(ctx *gin.Context, outCh chan<- string, errCh chan<- error) {
	parseStringParam(ctx, "categoryID", outCh, errCh)
}

// Parse query parameter to get integer variable.
func parseIntParam(ctx *gin.Context, name string, outCh chan<- int, errCh chan<- error) {
	defer close(outCh)
	str, ok := ctx.GetQuery(name)
	if !ok {
		errCh <- argus.NewError(errParamNotFound, nil)
		return
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		errCh <- argus.NewError(errParamIsNotNumber, err).
			AppendValue(name, str)
		return
	}

	outCh <- value
}

// Parse query parameter to get string variable.
func parseStringParam(ctx *gin.Context, name string, outCh chan<- string, errCh chan<- error) {
	defer close(outCh)
	value, ok := ctx.GetQuery(name)
	if !ok {
		errCh <- argus.NewError(errParamNotFound, nil).
			AppendValue("param", name)
		return
	}

	outCh <- value
}

func printMemory(mem runtime.MemStats) {
	runtime.ReadMemStats(&mem)
	fmt.Println(mem.Alloc, mem.TotalAlloc, mem.HeapAlloc, mem.HeapSys)
}
