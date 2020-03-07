package service

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/champon1020/argus"
)

func GenFlg(st interface{}, fieldNames ...string) (flg uint32) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for _, fn := range fieldNames {
		for j := 0; j < t.NumField(); j++ {
			if fn == t.Field(j).Name {
				flg |= 1 << j
			}
		}
	}
	return
}

func GenArgsSliceLogic(argsFlg uint32, st interface{}, isLimit bool) (args []interface{}) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if (1 << i & argsFlg) > 0 {
			args = append(args, v.Field(i).Interface())
		}
	}
	if isLimit {
		args = append(args, argus.GlobalConfig.Web.MaxViewArticleNum)
	}
	return
}

func GenArgsQuery(argsFlg uint32, st interface{}) (query string) {
	const initQuery = "WHERE "
	query = initQuery
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if (1 << i & argsFlg) > 0 {
			if query != initQuery {
				query += "AND "
			}
			query += ToSnakeCase(t.Field(i).Name) + "=" + "? "
		}
	}
	if query == initQuery {
		query = ""
	}
	return
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) (snake string) {
	snake = matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)
	return
}
