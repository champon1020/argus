package service_test

import (
	"testing"
	"time"

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
	if flg != actual {
		t.Errorf("mismatch flg: %v, actual: %v", actual, flg)
	}
}

func TestGenFlg_Id_Title(t *testing.T) {
	article := repository.Article{}

	flg := service.GenFlg(article, "Id", "Title")

	var actual uint32 = 3
	if flg != actual {
		t.Errorf("mismatch flg: %v, actual: %v", actual, flg)
	}
}

func TestGenArgsSliceLogic(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Title")
	st.Title = "test"
	args := service.GenArgsSliceLogic(argsFlg, st, false)

	if len(args) != 1 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual := "test"
	if args[0] != actual {
		t.Fatalf("value of args[0]: %v, actual: %v\n", args[0], actual)
	}
}

func TestGenArgsSliceLogic_Limit(t *testing.T) {
	var (
		argsFlg        uint32
		st             Hoge
		configurations argus.Configurations
	)
	configurations.New("dev")

	argsFlg = service.GenFlg(st, "Title")
	st.Title = "test"
	args := service.GenArgsSliceLogic(argsFlg, st, true)

	if len(args) != 2 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual1 := "test"
	if args[0] != actual1 {
		t.Fatalf("value of args[0]: %v, actual: %v\n", args[0], actual1)
	}

	config := argus.GlobalConfig
	actual2 := config.Web.MaxViewArticleNum
	if args[1] != actual2 {
		t.Fatalf("value of args[1]: %v, actual: %v\n", args[1], actual2)
	}
}

func TestGenArgsSliceLogic_Multi(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Id", "Title")
	st.Id = 1
	st.Title = "test"
	args := service.GenArgsSliceLogic(argsFlg, st, false)

	if len(args) != 2 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual1 := 1
	if args[0] != actual1 {
		t.Fatalf("value of args[0]: %v, actual: %v\n", args[0], actual1)
	}

	actual2 := "test"
	if args[1] != actual2 {
		t.Fatalf("value of args[1]: %v, actual: %v\n", args[0], actual2)
	}
}

func TestGenArgsQuery(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Title")
	args := service.GenArgsQuery(argsFlg, st)

	actual := "WHERE title=? "
	if args != actual {
		t.Fatalf("value of args: %v, actual: %v\n", args, actual)
	}
}

func TestGenArgsQuery_Multi(t *testing.T) {
	var (
		argsFlg uint32
		st      Hoge
	)
	argsFlg = service.GenFlg(st, "Title", "Date")
	args := service.GenArgsQuery(argsFlg, st)

	actual := "WHERE title=? AND date=? "
	if args != actual {
		t.Fatalf("value of args: %v, actual: %v\n", args, actual)
	}
}

func TestToSnakeCase(t *testing.T) {
	test := "TestTestTest012"
	actual := "test_test_test012"
	result := service.ToSnakeCase(test)
	if result != actual {
		t.Fatalf("result: %v, actual: %v\n", result, actual)
	}
}
