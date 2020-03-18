package service_test

import (
	"testing"
	"time"

	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/stretchr/testify/assert"
)

type Hoge struct {
	Id    int
	Title string
	Date  time.Time
}

func TestGenMask_Title(t *testing.T) {
	article := repo.Article{}
	fieldName := "Title"

	flg := service.GenMask(article, fieldName)

	var actual uint32 = 2
	assert.Equal(t, actual, flg)
}

func TestGenMask_Id_Title(t *testing.T) {
	article := repo.Article{}

	flg := service.GenMask(article, "Id", "Title")

	var actual uint32 = 3
	assert.Equal(t, actual, flg)
}

func TestGenArgsSlice(t *testing.T) {
	var (
		argsMask uint32
		st       Hoge
	)
	argsMask = service.GenMask(st, "Title")
	st.Title = "test"
	args := service.GenArgsSlice(argsMask, st, [2]int{})

	assert.Equal(t, 1, len(args))
	assert.Equal(t, "test", args[0])
}

func TestGenArgsSlice_Limit(t *testing.T) {
	var (
		argsMask uint32
		st       Hoge
	)

	argsMask = service.GenMask(st, "Title", "Limit")
	st.Title = "test"
	args := service.GenArgsSlice(argsMask, st, [2]int{1, 2})

	assert.Equal(t, 3, len(args))
	assert.Equal(t, "test", args[0])
	assert.Equal(t, 1, args[1])
	assert.Equal(t, 2, args[2])
}

func TestGenArgsSlice_Multi(t *testing.T) {
	var (
		argsMask uint32
		st       Hoge
	)
	argsMask = service.GenMask(st, "Id", "Title")
	st.Id = 1
	st.Title = "test"
	args := service.GenArgsSlice(argsMask, st, [2]int{})

	assert.Equal(t, 2, len(args))
	assert.Equal(t, 1, args[0])
	assert.Equal(t, "test", args[1])
}

func TestGenArgsQuery(t *testing.T) {
	var (
		argsMask uint32
		st       Hoge
	)
	argsMask = service.GenMask(st, "Title")
	args, limit := service.GenArgsQuery(argsMask, st)

	actual := "WHERE title=? "
	assert.Equal(t, actual, args+limit)
}

func TestGenArgsQuery_Multi(t *testing.T) {
	var (
		argsMask uint32
		st       Hoge
	)
	argsMask = service.GenMask(st, "Title", "Date", "Limit")
	args, limit := service.GenArgsQuery(argsMask, st)

	actual := "WHERE title=? AND date=? LIMIT ?,? "
	assert.Equal(t, actual, args+limit)
}

func TestToSnakeCase(t *testing.T) {
	test := "TestTestTest012"
	actual := "test_test_test012"
	result := service.ToSnakeCase(test)
	assert.Equal(t, actual, result)
}
