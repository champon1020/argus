package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/champon1020/argus"
)

var (
	Logger           = argus.Logger
	Errors           = &argus.Errors
	BasicError       = argus.NewError(argus.BasicError)
	IOReadError      = argus.NewError(argus.IOFailedReadError)
	IOMarshalError   = argus.NewError(argus.IOFailedMarshalError)
	IOUnmarshalError = argus.NewError(argus.IOFailedUnmarshalError)
	TimeParseError   = argus.NewError(argus.TimeFailedParseError)
)

func ReadBody(r io.Reader) (body []byte, err error) {
	if body, err = ioutil.ReadAll(r); err != nil {
		IOReadError.SetErr(err).AppendTo(Errors)
	}
	return
}

func ParseRequestBody(r *http.Request, reqBody *RequestBody) (err error) {
	var body []byte
	if body, err = ReadBody(r.Body); err != nil {
		return
	}
	if err = json.Unmarshal(body, &reqBody); err != nil {
		IOUnmarshalError.SetErr(err).AppendTo(Errors)
	}
	return
}

func ParseDraftRequestBody(r *http.Request, reqBody *DraftRequestBody) (err error) {
	var body []byte
	if body, err = ReadBody(r.Body); err != nil {
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
