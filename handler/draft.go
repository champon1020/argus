package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/champon1020/argus/repo"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

type DraftRequestArticle struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Categories  string `json:"categories"`
	ContentHash string `json:"contentHash"`
	ImageHash   string `json:"imageHash"`
}

type DraftRequestBody struct {
	Article    DraftRequestArticle `json:"article"`
	MdContents string              `json:"mdContents"`
}

type DraftInfoResp struct {
	Id          string `json:"id"`
	ContentHash string `json:"contentHash"`
	ImageHash   string `json:"imageHash"`
}

func DraftController(c *gin.Context) {
	_ = DraftHandler(c, repo.DraftCommand)
}

func DraftHandler(c *gin.Context, repoCmd repo.DraftCmd) (err error) {
	var (
		body     DraftRequestBody
		response string
	)

	if err = ParseDraftRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if body.Article.Id == "" {
		service.GenNewId(service.IdLen, &body.Article.Id)
	}

	fp := service.ResolveContentFilePath("drafts", body.Article.ContentHash)
	draft := repo.Draft{
		Id:          body.Article.Id,
		Title:       body.Article.Title,
		Categories:  body.Article.Categories,
		UpdateDate:  time.Now(),
		ContentHash: filepath.Base(fp),
		ImageHash:   body.Article.ImageHash,
	}

	mdFp := fp + "_md"
	if err = service.OutputFile(mdFp, []byte(body.MdContents)); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = repoCmd(*repo.GlobalMysql, draft); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		_ = service.DeleteFile(mdFp)
		return
	}

	res := DraftInfoResp{
		Id:          draft.Id,
		ContentHash: draft.ContentHash,
		ImageHash:   draft.ImageHash,
	}

	if response, err = ParseToJson(&res); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(c.Writer, response)
	return
}
