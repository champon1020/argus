package repository

import (
	"testing"
	"time"

	"github.com/champon1020/argus"
)

func TestGenArgsSlice(t *testing.T) {
	var (
		argsFlg uint32
		article Article
	)
	argsFlg = 1 << 1
	article.Title = "test"
	args := GenArgsSlice(argsFlg, article)

	if len(args) != 1 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual := "test"
	if args[0] != actual {
		t.Fatalf("value of args[0]: %v, actual: %v\n", args[0], actual)
	}
}

func TestGenArgsSliceIsLimit(t *testing.T) {
	var (
		argsFlg uint32
		article Article
		config  argus.Config
	)
	argsFlg = 1 << 1
	article.Title = "test"
	args := GenArgsSliceIsLimit(argsFlg, article, true)

	if len(args) != 2 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual1 := "test"
	if args[0] != actual1 {
		t.Fatalf("value of args[0]: %v, actual: %v\n", args[0], actual1)
	}

	config.Load()
	actual2 := config.Web.MaxViewArticleNum
	if args[1] != actual2 {
		t.Fatalf("value of args[1]: %v, actual: %v\n", args[1], actual2)
	}
}

func TestGenArgsSliceLogicTitle(t *testing.T) {
	var (
		argsFlg uint32
		article Article
	)
	argsFlg = 1 << 1
	article.Title = "test"
	args := GenArgsSliceLogic(argsFlg, article, false)

	if len(args) != 1 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual := "test"
	if args[0] != actual {
		t.Fatalf("value of args[0]: %v, actual: %v\n", args[0], actual)
	}
}

func TestGenArgsSliceLogicCreateDate(t *testing.T) {
	var (
		argsFlg uint32
		article Article
	)
	argsFlg = 1 << 3
	article.CreateDate, _ = time.Parse(time.RFC3339, "2006-01-02")
	args := GenArgsSliceLogic(argsFlg, article, false)

	if len(args) != 1 {
		t.Fatalf("length of args: %v\n", len(args))
	}

	actual, _ := time.Parse(time.RFC3339, "2006-01-02")
	if args[0] != actual {
		t.Fatalf("value of args[0]: %v, actual: %v\n", args[0], actual)
	}
}

func TestGenArgsQuery(t *testing.T) {
	var (
		argsFlg uint32
		article Article
	)
	argsFlg = 1 << 1
	args := GenArgsQuery(argsFlg, article)

	actual := "WHERE title=? "
	if args != actual {
		t.Fatalf("value of args: %v, actual: %v\n", args, actual)
	}
}

func TestGenArgsQueryTwo(t *testing.T) {
	var (
		argsFlg uint32
		article Article
	)
	argsFlg = 1<<1 | 1<<3
	args := GenArgsQuery(argsFlg, article)

	actual := "WHERE title=? AND create_date=? "
	if args != actual {
		t.Fatalf("value of args: %v, actual: %v\n", args, actual)
	}
}

func TestToSnakeCase(t *testing.T) {
	test := "TestTestTest012"
	actual := "test_test_test012"
	result := ToSnakeCase(test)
	if result != actual {
		t.Fatalf("result: %v, actual: %v\n", result, actual)
	}
}
