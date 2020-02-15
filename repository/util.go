package repository

import (
	"reflect"
	"regexp"
)

func GenArgsSlice(argsFlg uint32, st interface{}) (args []interface{}) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if 1 << i & argsFlg {
			args = append(args, v.Field(i))
		}
	}
	return
}

func GenArgsQuery(argsFlg uint32, st interface{}) (query string) {
	const initQuery = "WHERE "
	query = initQuery
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if 1 << i & argsFlg {
			if query != initQuery {
				query += "AND "
			}
			query += ToSnakeCase(v.Field(i).String()) + "=" + "? "
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
	return
}

func articleIdConverter(mysql MySQL, article *Article) {
	idList := GetEmptyMinId(mysql.DB, "article", 1)
	(*article).Id = idList[0]
}

func categoriesIdConverter(mysql MySQL, categories *[]Category) {
	idList := GetEmptyMinId(mysql.DB, "categories", len(*categories))
	for i := 0; i < len(*categories); i++ {
		(*categories)[i].Id = idList[i]
	}
}
