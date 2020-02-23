package repository

import (
	"reflect"
	"regexp"
	"strings"
)

func GenArgsSlice(argsFlg uint32, st interface{}) []interface{} {
	return GenArgsSliceLogic(argsFlg, st, false)
}

func GenArgsSliceIsLimit(argsFlg uint32, st interface{}, isLimit bool) []interface{} {
	return GenArgsSliceLogic(argsFlg, st, isLimit)
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
		args = append(args, config.Web.MaxViewArticleNum)
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

func ArticleIdConverter(mysql MySQL, article *Article) (err error) {
	var idList []int
	if idList, err = GetEmptyMinId(mysql.DB, "articles", 1); err != nil {
		return
	}

	(*article).Id = idList[0]
	return
}

func CategoriesIdConverter(mysql MySQL, categories *[]Category) (err error) {
	var idList []int
	if idList, err = GetEmptyMinId(mysql.DB, "categories", len(*categories)); err != nil {
		return
	}

	cur := 0
	for i := 0; i < len(*categories); i++ {
		if (*categories)[i].Id == -1 {
			(*categories)[i].Id = idList[cur]
			cur++
		}
	}
	return
}
