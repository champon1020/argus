package handler

import (
	"fmt"
	"net/http"
	"time"

	repo "github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

type DraftRequestArticle struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	Categories  []repo.Category `json:"categories"`
	ContentHash string          `json:"contentHash"`
	ImageHash   string          `json:"imageHash"`
}

type DraftRequestBody struct {
	Article  DraftRequestArticle `json:"article"`
	Contents string              `json:"contents"`
}

func DraftController(c *gin.Context) {
	DraftHandler(c, repo.DraftCommand)
}

func DraftHandler(c *gin.Context, repoCmd repo.DraftCmd) {
	var (
		body DraftRequestBody
		err  error
	)

	if err = ParseDraftRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fp := service.ResolveContentFilePath(body.Article.ContentHash, "drafts")
	draft := repo.Draft{
		Id:          body.Article.Id,
		Title:       body.Article.Title,
		Categories:  resolveToDraftCategories(body.Article.Categories),
		UpdateDate:  time.Now(),
		ContentHash: service.ConvertPathToFileName(fp),
		ImageHash:   body.Article.ImageHash,
	}

	if err = service.OutputFile(fp, body.Contents); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mysql := repo.GlobalMysql
	if err = repoCmd(mysql, draft); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		service.DeleteFile(fp)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
}

func resolveToDraftCategories(categories []repo.Category) string {
	res := ""
	for i, c := range categories {
		if i != 0 {
			res += "&"
		}
		res += c.Name
	}
	return res
}
