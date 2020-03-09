package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/champon1020/argus"
)

var (
	Logger           = argus.Logger
	Errors           = &argus.Errors
	IOReadError      = argus.NewError(argus.IOFailedReadError)
	IOMarshalError   = argus.NewError(argus.IOFailedMarshalError)
	IOUnmarshalError = argus.NewError(argus.IOFailedUnmarshalError)
)

func ParseRequestBody(r *http.Request, reqBody *RequestBody) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		IOReadError.SetErr(err).AppendTo(Errors)
		return
	}
	if err = json.Unmarshal(body, &reqBody); err != nil {
		IOUnmarshalError.SetErr(err).AppendTo(Errors)
		return
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
