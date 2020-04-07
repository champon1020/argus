package handler

import (
	"net/http"
	"time"

	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

type DraftRequestBody struct {
	Article struct {
		Id         string `json:"id"`
		Title      string `json:"title"`
		Categories string `json:"categories"`
		Content    string `json:"content"`
		ImageHash  string `json:"imageHash"`
	} `json:"article"`
}

type DraftInfoResp struct {
	Id        string `json:"id"`
	Content   string `json:"content"`
	ImageHash string `json:"imageHash"`
}

func DraftController(c *gin.Context) {
	_ = DraftHandler(c, repo.DraftCommand)
}

func DraftHandler(c *gin.Context, repoCmd repo.DraftCmd) (err error) {
	var body DraftRequestBody

	if err = ParseDraftRequestBody(c.Request, &body); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if body.Article.Id == "" {
		service.GenNewId(service.IdLen, &body.Article.Id)
	}

	draft := repo.Draft{
		Id:         body.Article.Id,
		Title:      body.Article.Title,
		Categories: body.Article.Categories,
		UpdateDate: time.Now(),
		Content:    body.Article.Content,
		ImageHash:  body.Article.ImageHash,
	}

	if err = repoCmd(*repo.GlobalMysql, draft); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := DraftInfoResp{
		Id:        draft.Id,
		Content:   draft.Content,
		ImageHash: draft.ImageHash,
	}

	c.JSON(http.StatusOK, res)
	return
}
