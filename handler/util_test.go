package handler

import (
	"bytes"
	"net/http"
	"testing"
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
