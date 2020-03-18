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

// map[<Struct Field Name>]<Operation>
type ArgsOpeMap map[string]Ope

// Database query information.
type QueryOption struct {
	Args   []interface{}
	Aom    ArgsOpeMap
	Limit  int
	Offset int
	Order  string
	Desc   bool
}

func (op *QueryOption) BuildArgs() {
	if op.Limit != 0 {
		op.Args = append(op.Args, op.Offset, op.Limit)
	}
}

var DefaultOption = &QueryOption{
	Args:   []interface{}{},
	Aom:    map[string]Ope{},
	Limit:  0,
	Offset: 0,
	Order:  "",
	Desc:   false,
}

// Generate arguments query used in database query.
// Return values is query(query of following 'WHERE')
func GenArgsQuery(option QueryOption) (query string) {
	for k, v := range option.Aom {
		if query == "" {
			query += "WHERE "
		} else {
			query += "AND "
		}
		query += ToSnakeCase(k) + v.toString() + "? "
	}
	if option.Order != "" {
		query += "ORDER BY " + option.Order + " "
		if option.Desc {
			query += "DESC "
		}
	}
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
