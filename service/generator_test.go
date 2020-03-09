package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/service"
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

func TestGenArgsSliceLogic(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Title")
	st.Title = "test"
	args := service.GenArgsSliceLogic(argsFlg, st)

	if len(args) != 1 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual := "test"
	assert.Equal(t, actual, args[0])
}

func TestGenArgsSliceLogic_Limit(t *testing.T) {
	var (
		argsFlg        uint32
		st             Hoge
		configurations argus.Configurations
	)
	configurations.New("dev")

	argsFlg = service.GenFlg(st, "Title", "Limit")
	st.Title = "test"
	args := service.GenArgsSliceLogic(argsFlg, st)

	if len(args) != 2 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual1 := "test"
	assert.Equal(t, actual1, args[0])

	config := argus.GlobalConfig
	actual2 := config.Web.MaxViewArticleNum
	assert.Equal(t, actual2, args[1])
}

func TestGenArgsSliceLogic_Multi(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Id", "Title")
	st.Id = 1
	st.Title = "test"
	args := service.GenArgsSliceLogic(argsFlg, st)

	if len(args) != 2 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual1 := 1
	assert.Equal(t, actual1, args[0])

	actual2 := "test"
	assert.Equal(t, actual2, args[1])
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

	actual := "WHERE title=? AND date=? LIMIT ? "
	assert.Equal(t, actual, args+limit)
}

func TestToSnakeCase(t *testing.T) {
	test := "TestTestTest012"
	actual := "test_test_test012"
	result := service.ToSnakeCase(test)
	assert.Equal(t, actual, result)
}
