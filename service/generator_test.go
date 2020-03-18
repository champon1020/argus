package service_test

import (
	"testing"

	"github.com/champon1020/argus/service"
	"github.com/stretchr/testify/assert"
)

type Hoge struct {
	Id    int
	Title string
}

func TestQueryOption_BuildArgs(t *testing.T) {
	option := &service.QueryOption{
		Args:   []interface{}{"test"},
		Aom:    service.ArgsOpeMap{"Title": service.Eq},
		Limit:  2,
		Offset: 3,
		Order:  "",
		Desc:   false,
	}

	option.BuildArgs()
	assert.Equal(t, 3, len(option.Args))
	assert.Equal(t, "test", option.Args[0])
	assert.Equal(t, 3, option.Args[1])
	assert.Equal(t, 2, option.Args[2])
}

func TestGenArgsQuery(t *testing.T) {
	option := &service.QueryOption{
		Args:   []interface{}{"test"},
		Aom:    service.ArgsOpeMap{"Title": service.Eq},
		Limit:  0,
		Offset: 0,
		Order:  "",
		Desc:   false,
	}

	query := service.GenArgsQuery(*option)
	actual := "WHERE title=? "
	assert.Equal(t, actual, query)
}

func TestGenArgsQuery_Multi(t *testing.T) {
	option := &service.QueryOption{
		Args:   []interface{}{1, "test"},
		Aom:    service.ArgsOpeMap{"Id": service.Ge, "Title": service.Eq},
		Limit:  3,
		Offset: 2,
		Order:  "create_date",
		Desc:   true,
	}

	query := service.GenArgsQuery(*option)
	actual := "WHERE id>=? AND title=? ORDER BY create_date DESC LIMIT ?,? "
	assert.Equal(t, actual, query)
}

func TestToSnakeCase(t *testing.T) {
	test := "TestTestTest012"
	actual := "test_test_test012"
	result := service.ToSnakeCase(test)
	assert.Equal(t, actual, result)
}
