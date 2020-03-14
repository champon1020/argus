package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/champon1020/argus"
	repo "github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	loc, _   = time.LoadLocation("Asia/Tokyo")
	testTime = time.Date(2020, 3, 9, 0, 0, 0, 0, loc)
)

func TestMain(m *testing.M) {
	argus.GlobalConfig = argus.NewConfig("dev")
	repo.GlobalMysql = repo.NewMysql()

	ret := m.Run()
	os.Exit(ret)
}

func TestFindArticleHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/find/article/list",
		nil)

	repoCmdMock := func(_ repo.MySQL, _ repo.Article, _ uint32) (articles []repo.Article, _ error) {
		articles = append(articles, repo.Article{
			Id:    1,
			Title: "test",
			Categories: []repo.Category{
				{1, "c1"},
			},
			CreateDate:  testTime,
			UpdateDate:  testTime,
			ContentHash: "0123456789",
			ImageHash:   "9876543210",
			Private:     false,
		})
		return
	}

	expectedBody := `{
	"articles": [
		{
			"id": 1,
			"title": "test",
			"categories": [
				{
					"id": 1,
					"name": "c1"
				}
			],
			"createDate": "2020-03-09T00:00:00+09:00",
			"updateDate": "2020-03-09T00:00:00+09:00",
			"contentHash": "0123456789",
			"imageHash": "9876543210",
			"private": false
		}
	]
}`

	if err := FindArticleHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}

func TestFindArticleByIdHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/find/article/list/id?id=2",
		nil)

	repoCmdMock := func(_ repo.MySQL, _ repo.Article, _ uint32) (articles []repo.Article, _ error) {
		articles = append(articles, repo.Article{
			Id:    2,
			Title: "test",
			Categories: []repo.Category{
				{1, "c1"},
			},
			CreateDate:  testTime,
			UpdateDate:  testTime,
			ContentHash: "0123456789",
			ImageHash:   "9876543210",
			Private:     false,
		})
		return
	}

	expectedBody := `{
	"articles": [
		{
			"id": 2,
			"title": "test",
			"categories": [
				{
					"id": 1,
					"name": "c1"
				}
			],
			"createDate": "2020-03-09T00:00:00+09:00",
			"updateDate": "2020-03-09T00:00:00+09:00",
			"contentHash": "0123456789",
			"imageHash": "9876543210",
			"private": false
		}
	]
}`

	if err := FindArticleByIdHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}

func TestFindArticleByTitleHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/find/article/list/title?title=test",
		nil)

	repoCmdMock := func(_ repo.MySQL, a repo.Article, _ uint32) (articles []repo.Article, err error) {
		if a.Title != "test" {
			err = errors.New("query parameter is not valid")
			return
		}
		articles = append(articles, repo.Article{
			Id:    1,
			Title: "test",
			Categories: []repo.Category{
				{1, "c1"},
			},
			CreateDate:  testTime,
			UpdateDate:  testTime,
			ContentHash: "0123456789",
			ImageHash:   "9876543210",
			Private:     false,
		})
		return
	}

	expectedBody := `{
	"articles": [
		{
			"id": 1,
			"title": "test",
			"categories": [
				{
					"id": 1,
					"name": "c1"
				}
			],
			"createDate": "2020-03-09T00:00:00+09:00",
			"updateDate": "2020-03-09T00:00:00+09:00",
			"contentHash": "0123456789",
			"imageHash": "9876543210",
			"private": false
		}
	]
}`

	if err := FindArticleByTitleHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}

func TestFindArticleByCreateDateHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/find/article/list/create-date?createDate=2020-03-09T00:00:00Z",
		nil)

	repoCmdMock := func(_ repo.MySQL, a repo.Article, _ uint32) (articles []repo.Article, err error) {
		if testTime.Equal(a.CreateDate) {
			err = errors.New("query parameter is not valid")
			return
		}
		articles = append(articles, repo.Article{
			Id:    1,
			Title: "test",
			Categories: []repo.Category{
				{1, "c1"},
			},
			CreateDate:  testTime,
			UpdateDate:  testTime,
			ContentHash: "0123456789",
			ImageHash:   "9876543210",
			Private:     false,
		})
		return
	}

	expectedBody := `{
	"articles": [
		{
			"id": 1,
			"title": "test",
			"categories": [
				{
					"id": 1,
					"name": "c1"
				}
			],
			"createDate": "2020-03-09T00:00:00+09:00",
			"updateDate": "2020-03-09T00:00:00+09:00",
			"contentHash": "0123456789",
			"imageHash": "9876543210",
			"private": false
		}
	]
}`

	if err := FindArticleByCreateDateHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}

func TestFindArticleByCategoryHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/find/article/list/category?category=c1&category=c2",
		nil)

	repoCmdMock := func(_ repo.MySQL, caNames []string, _ uint32) (articles []repo.Article, err error) {
		if len(caNames) != 2 {
			err = errors.New("category names length is not valid")
		}
		if caNames[0] != "c1" || caNames[1] != "c2" {
			err = errors.New("category names are not valid")
		}
		articles = append(articles, repo.Article{
			Id:    1,
			Title: "test",
			Categories: []repo.Category{
				{1, "c1"},
				{2, "c2"},
			},
			CreateDate:  testTime,
			UpdateDate:  testTime,
			ContentHash: "0123456789",
			ImageHash:   "9876543210",
			Private:     false,
		})
		return
	}

	expectedBody := `{
	"articles": [
		{
			"id": 1,
			"title": "test",
			"categories": [
				{
					"id": 1,
					"name": "c1"
				},
				{
					"id": 2,
					"name": "c2"
				}
			],
			"createDate": "2020-03-09T00:00:00+09:00",
			"updateDate": "2020-03-09T00:00:00+09:00",
			"contentHash": "0123456789",
			"imageHash": "9876543210",
			"private": false
		}
	]
}`

	if err := FindArticleByCategoryHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}

func TestFindCategoryHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/find/category/list",
		nil)

	repoCmdMock := func(_ repo.MySQL, _ repo.Category, _ uint32) (categories []repo.CategoryResponse, _ error) {
		categories = append(categories, repo.CategoryResponse{
			Id:         1,
			Name:       "c1",
			ArticleNum: 3,
		})
		return
	}

	expectedBody := `{
	"categories": [
		{
			"id": 1,
			"name": "c1",
			"articleNum": 3
		}
	]
}`

	if err := FindCategoryHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}

func TestFindDraftHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/find/draft",
		nil)

	repoCmdMock := func(_ repo.MySQL, _ repo.Draft, _ uint32) (drafts []repo.Draft, _ error) {
		drafts = append(drafts, repo.Draft{
			Id:          1,
			Title:       "test",
			Categories:  "c1&c2",
			UpdateDate:  testTime,
			ContentHash: "0123456789",
			ImageHash:   "9876543210",
		})
		return
	}

	expectedBody := `{
	"drafts": [
		{
			"id": 1,
			"title": "test",
			"categories": "c1\u0026c2",
			"updateDate": "2020-03-09T00:00:00+09:00",
			"contentHash": "0123456789",
			"imageHash": "9876543210"
		}
	]
}`

	if err := FindDraftHandler(ctx, repoCmdMock); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}

func TestFindImageHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(
		"GET",
		"/api/find/image/list",
		nil)

	expectedBody := `{
	"images": [
		"image_test1.png",
		"image_test2.jpg"
	],
	"next": false
}`

	if err := FindImageHandler(ctx); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("error happend in handler")
	}

	res := w.Result()
	assert.Equal(t, res.StatusCode, 200)

	var buf bytes.Buffer
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Indent(&buf, body, "", "	"); err != nil {
		argus.StdLogger.ErrorLog(*Errors)
		t.Fatalf("Unable to indent json string\n")
		return
	}
	assert.Equal(t, expectedBody, buf.String())
	buf.Reset()
}
