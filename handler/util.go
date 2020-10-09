package handler

import (
	"errors"
	"fmt"
	"net/http"
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
func ParsePage(ctx *gin.Context, outc chan<- int, errc chan<- error) {
	parseIntParam(ctx, "p", outc, errc)
}

// ParseNum parses query parameter to get title string.
func ParseNum(ctx *gin.Context, outc chan<- int, errc chan<- error) {
	parseIntParam(ctx, "num", outc, errc)
}

// ParseSortedID parses query parameter to get title string.
func ParseSortedID(ctx *gin.Context, outc chan<- int, errc chan<- error) {
	parseIntParam(ctx, "sortedID", outc, errc)
}

// ParseTitle parses query parameter to get title string.
func ParseTitle(ctx *gin.Context, outc chan<- string, errc chan<- error) {
	parseStringParam(ctx, "title", outc, errc)
}

// ParseID parses query parameter to get string id.
func ParseID(ctx *gin.Context, outc chan<- string, errc chan<- error) {
	parseStringParam(ctx, "id", outc, errc)
}

// ParseCategoryID parses query parameter to get category id.
func ParseCategoryID(ctx *gin.Context, outc chan<- string, errc chan<- error) {
	parseStringParam(ctx, "categoryID", outc, errc)
}

// Parse query parameter to get integer variable.
func parseIntParam(ctx *gin.Context, name string, outc chan<- int, errc chan<- error) {
	defer close(outc)
	str, ok := ctx.GetQuery(name)
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		err := argus.NewError(errParamNotFound, nil)
		errc <- err
		return
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		err := argus.NewError(errParamIsNotNumber, err).
			AppendValue(name, str)
		errc <- err
		return
	}

	outc <- value
}

// Parse query parameter to get string variable.
func parseStringParam(ctx *gin.Context, name string, outc chan<- string, errc chan<- error) {
	defer close(outc)
	value, ok := ctx.GetQuery(name)
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		err := argus.NewError(errParamNotFound, nil).
			AppendValue("param", name)
		errc <- err
		return
	}

	outc <- value
}

func printMemory(mem runtime.MemStats) {
	runtime.ReadMemStats(&mem)
	fmt.Println(mem.Alloc, mem.TotalAlloc, mem.HeapAlloc, mem.HeapSys)
}
