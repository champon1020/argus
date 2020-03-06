package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/champon1020/argus/repository"

	"github.com/champon1020/argus"
)

type RequestBody struct {
	Article  repository.Article `json:"article"`
	Contents string             `json:"contents"`
}

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

func GenFlg(st interface{}, fieldNames ...string) (flg uint32) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for _, fn := range fieldNames {
		for j := 0; j < t.NumField(); j++ {
			if fn == t.Field(j).Name {
				flg |= 1 << (j + 1)
			}
		}
	}
	return
}
