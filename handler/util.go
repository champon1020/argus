package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/champon1020/argus/service"

	"github.com/champon1020/argus"
)

var (
	Errors           = &argus.Errors
	BasicError       = argus.NewError(argus.BasicError)
	IOMarshalError   = argus.NewError(argus.IOFailedMarshalError)
	IOUnmarshalError = argus.NewError(argus.IOFailedUnmarshalError)
	TimeParseError   = argus.NewError(argus.TimeFailedParseError)
)

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

func ParseToJson(st interface{}) (res string, err error) {
	var bytes []byte
	if bytes, err = json.Marshal(&st); err != nil {
		IOMarshalError.SetErr(err).AppendTo(Errors)
		return
	}
	res = string(bytes)
	return
}
