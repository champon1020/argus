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
			"id": "TEST_ID",
			"title": "test",
			"categories": [
				{
					"id": "TEST_CA_ID",
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

	assert.Equal(t, "TEST_ID", body.Article.Id)
	assert.Equal(t, "test", body.Article.Title)
	assert.Equal(t, 1, len(body.Article.Categories))
	assert.Equal(t, "TEST_CA_ID", body.Article.Categories[0].Id)
	assert.Equal(t, "test_test", body.Article.Categories[0].Name)
	assert.Equal(t, "http://localhost:2000/", body.Article.ContentHash)
	assert.Equal(t, "http://localhost:1000/", body.Article.ImageHash)
	assert.Equal(t, false, body.Article.Private)
	assert.Equal(t, "<div>ok</div>", body.Contents)
}

func TestParseOffsetLimit(t *testing.T) {
	p := 1
	aNum := 10
	ol := ParseOffsetLimit(p, aNum)
	assert.Equal(t, 0, ol[0])
	assert.Equal(t, aNum, ol[1])

	p = 2
	aNum = 10
	ol = ParseOffsetLimit(p, aNum)
	assert.Equal(t, 10, ol[0])
	assert.Equal(t, aNum, ol[1])
}

func TestMax(t *testing.T) {
	a := 2
	b := 3
	c := Max(a, b)
	assert.Equal(t, 3, c)
}

func TestGetMaxPage(t *testing.T) {
	mxView := 3
	num := 5
	c := GetMaxPage(num, mxView)
	assert.Equal(t, 2, c)

	mxView = 4
	num = 12
	c = GetMaxPage(num, mxView)
	assert.Equal(t, 3, c)
}

func TestParseDraftRequestBody(t *testing.T) {
	var body DraftRequestBody

	requestJson := `{
		"article": {
			"id": "TEST_ID",
			"title": "test",
			"categories": [
				{
					"id": "TEST_CA_ID",
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

	assert.Equal(t, "TEST_ID", body.Article.Id)
	assert.Equal(t, "test", body.Article.Title)
	assert.Equal(t, 1, len(body.Article.Categories))
	assert.Equal(t, "TEST_CA_ID", body.Article.Categories[0].Id)
	assert.Equal(t, "test_test", body.Article.Categories[0].Name)
	assert.Equal(t, "http://localhost:2000/", body.Article.ContentHash)
	assert.Equal(t, "http://localhost:1000/", body.Article.ImageHash)
	assert.Equal(t, "<div>ok</div>", body.Contents)
}
