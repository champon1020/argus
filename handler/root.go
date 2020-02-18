package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repository"
)

var (
	mysql  repository.MySQL
	config argus.Config
	DBNAME = "argus"
)

func init() {
	config.Load()
	if err := mysql.Connect(config.DevDb, DBNAME); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request, method string) (body RequestBody, err error) {
	return rootHandler(w, r, method, "application/json")
}

func rootHandler(w http.ResponseWriter, r *http.Request, method string, contentType string) (body RequestBody, err error) {
	isErr := Validation(&w, r, method, contentType)
	if isErr {
		err = errors.New("unable to pass validation")
		return
	}

	if err = ParseRequestBody(&w, r, &body); err != nil {
		return
	}

	return
}
