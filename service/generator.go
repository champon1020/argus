package service

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Mathematical operation of database query.
type Ope string

func (o *Ope) toString() string {
	return string(*o)
}

const (
	Ne Ope = "!=" // Not Equal
	Eq Ope = "="  // Equal
	Gt Ope = ">"  // Greater Than
	Lt Ope = "<"  // Less Than
	Ge Ope = ">=" // Greater or Equal
	Le Ope = "<=" // Less or Equal
)

type QueryArgs struct {
	Value interface{}
	Name  string
	Ope   Ope
}

// Database query information.
type QueryOption struct {
	Args   []*QueryArgs
	Limit  int
	Offset int
	Order  string
	Desc   bool
}

var DefaultOption = &QueryOption{
	Args:   []*QueryArgs{},
	Limit:  0,
	Offset: 0,
	Order:  "",
	Desc:   false,
}

func GenArgsSlice(option QueryOption) (args []interface{}) {
	for _, a := range option.Args {
		args = append(args, a.Value)
	}
	if option.Limit != 0 {
		args = append(args, option.Offset, option.Limit)
	}
	return
}

// Generate arguments query used in database query.
// Return values is query(query of following 'WHERE')
func GenArgsQuery(option QueryOption) string {
	return GenWhereQuery(option) + GenOrderQuery(option) + GenLimitQuery(option)
}

func GenWhereQuery(option QueryOption) (query string) {
	for _, a := range option.Args {
		if query == "" {
			query += "WHERE "
		} else {
			query += "AND "
		}
		query += ToSnakeCase(a.Name) + a.Ope.toString() + "? "
	}
	return
}

func GenOrderQuery(option QueryOption) (query string) {
	if option.Order != "" {
		query += "ORDER BY " + option.Order + " "
		if option.Desc {
			query += "DESC "
		}
	}
	return
}

func GenLimitQuery(option QueryOption) (query string) {
	if option.Limit != 0 {
		query += "LIMIT ?,? "
	}
	return
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const IdLen = 50

func GenNewId(length int, id *string) {
	if *id == "TEST_ID" || *id == "TEST_CA_ID" {
		return
	}
	var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randSeed.Intn(len(charset))]
	}
	*id = string(b)
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// Convert camel case string to snake case.
func ToSnakeCase(str string) (snake string) {
	str = strings.Split(str, "#")[0]
	snake = matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)
	return
}
