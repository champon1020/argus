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
	logger         argus.Logger
	mysql          repository.MySQL
	configurations argus.Configurations
	config         argus.Config
	DBNAME         = "argus"
)

func init() {
	logger.NewLogger("[handler]")
	config = configurations.Load()
	if err := mysql.Connect(config.Db); err != nil {
		log.Fatalf("%v\n", err)
	}
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
