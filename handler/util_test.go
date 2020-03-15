package handler

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, 2, body.Article.Id)
	assert.Equal(t, "test", body.Article.Title)
	assert.Equal(t, 1, len(body.Article.Categories))
	assert.Equal(t, 1, body.Article.Categories[0].Id)
	assert.Equal(t, "http://localhost:2000/", body.Article.ContentHash)
	assert.Equal(t, "http://localhost:1000/", body.Article.ImageHash)
	assert.Equal(t, false, body.Article.Private)
	assert.Equal(t, "<div>ok</div>", body.Contents)
}

func TestMax(t *testing.T) {
	a := 2
	b := 3
	c := Max(a, b)
	assert.Equal(t, 3, c)
}

func TestParseDraftRequestBody(t *testing.T) {
	var body DraftRequestBody

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
			"imageHash": "http://localhost:1000/"
		},
		"contents": "<div>ok</div>"
	}`

	r, _ := http.NewRequest("POST", "", bytes.NewBuffer([]byte(requestJson)))
	if err := ParseDraftRequestBody(r, &body); err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, body.Article.Id)
	assert.Equal(t, "test", body.Article.Title)
	assert.Equal(t, 1, len(body.Article.Categories))
	assert.Equal(t, 1, body.Article.Categories[0].Id)
	assert.Equal(t, "http://localhost:2000/", body.Article.ContentHash)
	assert.Equal(t, "http://localhost:1000/", body.Article.ImageHash)
	assert.Equal(t, "<div>ok</div>", body.Contents)
}
