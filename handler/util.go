package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/champon1020/argus/repository"

	"github.com/champon1020/argus"
)

type RequestBody struct {
	Article repository.Article `json:"article"`
}

var (
	logger argus.Logger
	mysql  repository.MySQL
	config argus.Config
	DBNAME = "argus"
)

func init() {
	logger.NewLogger("[handler]")
	config.Load()
	if err := mysql.Connect(config.DevDb, DBNAME); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func Validation(w *http.ResponseWriter, r *http.Request, method string, contentType string) (isErr bool) {
	isErr = true
	if r.Method != method {
		(*w).WriteHeader(http.StatusBadRequest)
		logger.Printf("This is not %s request\n", method)
		return
	}

	if r.Header.Get("Content-Type") != contentType {
		(*w).WriteHeader(http.StatusBadRequest)
		logger.Printf("Content-Type is not %s\n", contentType)
		return
	}

	if l := r.ContentLength; l == 0 {
		(*w).WriteHeader(http.StatusInternalServerError)
		logger.Println("Content-Length is 0")
		return
	}

	isErr = false
	return
}

func ParseRequestBody(w *http.ResponseWriter, r *http.Request, entity *RequestBody) error {
	var (
		body []byte
		err  error
	)
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		logger.Println("Unable to read request body")
		return err
	}

	if err = json.Unmarshal(body, &entity); err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		logger.Println("Json format is not comfortable")
		return err
	}
	return nil
}

func GenFlg(st interface{}, fieldNames ...string) (flg int) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for _, fn := range fieldNames {
		for j := 0; j < t.NumField(); j++ {
			if fn == v.Field(j).String() {
				flg |= 1 << j
			}
		}
	}
	return
}
