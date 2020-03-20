package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

var (
	Errors           = &argus.Errors
	BasicError       = argus.NewError(argus.BasicError)
	IOMarshalError   = argus.NewError(argus.IOFailedMarshalError)
	IOUnmarshalError = argus.NewError(argus.IOFailedUnmarshalError)
	TimeParseError   = argus.NewError(argus.TimeFailedParseError)
)

// Parse page from context of gin framework.
func ParsePage(c *gin.Context) (p int, err error) {
	p = 1
	if pp, ok := c.GetQuery("p"); ok {
		if p, err = strconv.Atoi(pp); err != nil {
			BasicError.SetErr(err).AppendTo(Errors)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	return
}

// Get LimitOffset object from p (page num).
func ParseOffsetLimit(p int, maxView int) (ol repo.OffsetLimit) {
	ol[1] = maxView
	ol[0] = (p - 1) * ol[1]
	return
}

// Get max of a and b;
func Max(a int, b int) (r int) {
	if a >= b {
		r = a
	} else {
		r = b
	}
	return
}

// Get total page of view
func GetMaxPage(articlesNum int, maxView int) int {
	return (maxView + articlesNum - 1) / maxView
}

// Parse request object from http.Request.
func ParseRequestBody(r *http.Request, reqBody *RequestBody) (err error) {
	var body []byte
	if body, err = service.ReadBody(r.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, &reqBody); err != nil {
		IOUnmarshalError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Parse draft request object from http.Request.
func ParseDraftRequestBody(r *http.Request, reqBody *DraftRequestBody) (err error) {
	var body []byte
	if body, err = service.ReadBody(r.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, &reqBody); err != nil {
		IOUnmarshalError.SetErr(err).AppendTo(Errors)
	}
	return
}

// Parse json string from object.
func ParseToJson(st interface{}) (res string, err error) {
	var bytes []byte
	if bytes, err = json.Marshal(&st); err != nil {
		IOMarshalError.SetErr(err).AppendTo(Errors)
		return
	}
	res = string(bytes)
	return
}
