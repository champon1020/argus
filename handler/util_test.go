package handler

import (
	"bytes"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/champon1020/argus/service"

	"github.com/champon1020/argus/repository"
)

func TestParseRequestBody(t *testing.T) {
	var body RequestBody

	requestJson := `{
		"article": {
			"id": 2,
			"title": "test",
			"categories": [
				{
					"id": 1,
					"name": "test_test"
				}
			],
			"createDate": "2018-01-02T00:00:00+09:00",
			"updateDate": "2018-01-03T00:00:00+09:00",
			"contentHash": "http://localhost:2000/",
			"imageHash": "http://localhost:1000/",
			"private": false
		},
		"contents": "<div>ok</div>"
	}`

	r, _ := http.NewRequest("POST", "", bytes.NewBuffer([]byte(requestJson)))

	if err := ParseRequestBody(r, &body); err != nil {
		t.Error(err)
	}

	if body.Article.Id != 2 {
		t.Errorf("mismatch id: %v, but actual: %v", 2, body.Article.Id)
	}
	if body.Article.Title != "test" {
		t.Errorf("mismatch title: %v, but actual: %v", "test", body.Article.Title)
	}
	if len(body.Article.Categories) == 0 {
		t.Errorf("categories is empty")
	}
	if body.Article.Categories[0].Id != 1 {
		t.Errorf("mismatch category id: %v, but actual: %v", 1, body.Article.Categories[0].Id)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	createDate := time.Date(2018, 1, 2, 0, 0, 0, 0, loc)
	if !body.Article.CreateDate.Equal(createDate) {
		t.Errorf("mismatch create date: %v, but actual: %v", createDate, body.Article.CreateDate)
	}

	updateDate := time.Date(2018, 1, 3, 0, 0, 0, 0, loc)
	if !body.Article.UpdateDate.Equal(updateDate) {
		t.Errorf("mismatch update date: %v, but actual: %v", updateDate, body.Article.UpdateDate)
	}

	if body.Article.ContentHash != "http://localhost:2000/" {
		t.Errorf("mismatch content url: %v, but actual: %v", "http://localhost:2000/", body.Article.ContentHash)
	}
	if body.Article.ImageHash != "http://localhost:1000/" {
		t.Errorf("mismatch image url: %v, but actual: %v", "http://localhost:1000/", body.Article.ImageHash)
	}
	if body.Article.Private != false {
		t.Errorf("mismatch category id: %v, but actual: %v", false, body.Article.Private)
	}

	if body.Contents != "<div>ok</div>" {
		t.Errorf("mismatch contents: %v, but actual: %v", "<div>ok</div>", body.Contents)
	}
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

func TestResolveContentHash(t *testing.T) {
	contentHash := "article"
	fn := ResolveContentHash(contentHash)
	assert.Equal(t, "article", fn)

	contentHash = ""
	fn = ResolveContentHash(contentHash)
	today := time.Now()
	tStr := today.Format("20060102150405")
	assert.Equal(t, string([]rune(tStr)[:6]), string([]rune(fn)[:6]))
}

func TestResolveContentFilePath(t *testing.T) {
	contentHash := "test"
	dirName := "articles"
	path := ResolveContentFilePath(contentHash, dirName)
	assert.Equal(t, os.Getenv("ARGUS_ARTICLE")+"/articles/test", path)
}

func TestConvertPathToFileName(t *testing.T) {
	path := "/this/is/test/article"
	fn := ConvertPathToFileName(path)
	assert.Equal(t, "article", fn)
}
