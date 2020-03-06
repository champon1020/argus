package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/champon1020/argus"
)

var logger argus.Logger

func init() {
	logger.NewLogger("[handler]")
}

func ParseRequestBody(r *http.Request, reqBody *RequestBody) error {
	var (
		body []byte
		err  error
	)

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		logger.Println("Unable to read request body")
		return err
	}

	if err = json.Unmarshal(body, &reqBody); err != nil {
		logger.Println("Json format is not comfortable")
		return err
	}
	return nil
}

func ParseToJson(st interface{}) (res string, err error) {
	var bytes []byte
	if bytes, err = json.Marshal(&st); err != nil {
		logger.ErrorPrintf(err)
		return
	}
	res = string(bytes)
	return
}
