package service_test

import (
	"testing"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/service"
	"github.com/stretchr/testify/assert"
)

type Hoge struct {
	Id    int
	Title string
	Date  time.Time
}

func TestGenFlg_Title(t *testing.T) {
	article := repository.Article{}
	fieldName := "Title"

	flg := service.GenFlg(article, fieldName)

	var actual uint32 = 2
	assert.Equal(t, actual, flg)
}

func TestGenFlg_Id_Title(t *testing.T) {
	article := repository.Article{}

	flg := service.GenFlg(article, "Id", "Title")

	var actual uint32 = 3
	assert.Equal(t, actual, flg)
}

func TestGenArgsSlice(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Title")
	st.Title = "test"
	args := service.GenArgsSlice(argsFlg, st, [2]int{})

	assert.Equal(t, 1, len(args))
	assert.Equal(t, "test", args[0])
}

func TestGenArgsSlice_Limit(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argus.GlobalConfig = argus.NewConfig("dev")

	argsFlg = service.GenFlg(st, "Title", "Limit")
	st.Title = "test"
	args := service.GenArgsSlice(argsFlg, st, [2]int{1, 2})

	assert.Equal(t, 3, len(args))
	assert.Equal(t, "test", args[0])
	assert.Equal(t, 1, args[1])
	assert.Equal(t, 2, args[2])
}

func TestGenArgsSlice_Multi(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Id", "Title")
	st.Id = 1
	st.Title = "test"
	args := service.GenArgsSlice(argsFlg, st, [2]int{})

	assert.Equal(t, 2, len(args))
	assert.Equal(t, 1, args[0])
	assert.Equal(t, "test", args[1])
}

func TestGenArgsQuery(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Title")
	args, limit := service.GenArgsQuery(argsFlg, st)

	actual := "WHERE title=? "
	assert.Equal(t, actual, args+limit)
}

func TestGenArgsQuery_Multi(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Title", "Date", "Limit")
	args, limit := service.GenArgsQuery(argsFlg, st)

	actual := "WHERE title=? AND date=? LIMIT ?,? "
	assert.Equal(t, actual, args+limit)
}

func TestToSnakeCase(t *testing.T) {
	test := "TestTestTest012"
	actual := "test_test_test012"
	result := service.ToSnakeCase(test)
	assert.Equal(t, actual, result)
}
