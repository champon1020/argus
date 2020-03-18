package service

import (
	"reflect"
	"regexp"
	"strings"
)

// Generate flag which determines arguments of database query.
// For example, if you want a query like '... WHERE title=?',
// you should set 'Title' of string to fieldNames.
// But selected struct 'st' must have a field named 'Title'.
func GenMask(st interface{}, fieldNames ...string) (mask uint32) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for _, fn := range fieldNames {
		if fn == "Limit" {
			mask |= 1 << 31
			continue
		}
		for j := 0; j < t.NumField(); j++ {
			if fn == t.Field(j).Name {
				mask |= 1 << j
			}
		}
	}
	return
}

// Generate arguments slice used in database query.
// If (flg & 1 << 31) > 0, limit query is turned on.
func GenArgsSlice(argsFlg uint32, st interface{}, ol ...[2]int) (args []interface{}) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if (1 << i & argsFlg) > 0 {
			args = append(args, v.Field(i).Interface())
		}
	}
	if (1<<31&argsFlg) > 0 && len(ol) >= 1 {
		args = append(args, ol[0][0], ol[0][1])
	}
	return
}

// Generate arguments query used in database query.
// If (flg & 1 << 31) > 0, limit query is turned on.
// Return values is query(query of following 'WHERE')
// and limit(limit query 'LIMIT ?,?').
func GenArgsQuery(argsFlg uint32, st interface{}) (query string, limit string) {
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
	if (1 << 31 & argsFlg) > 0 {
		limit = "LIMIT ?,? "
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

// Convert camel case string to snake case.
func ToSnakeCase(str string) (snake string) {
	snake = matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)
	return
}
