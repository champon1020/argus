package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/champon1020/argus/repository"

	"github.com/gin-gonic/gin"
)

type DraftRequestArticle struct {
	Id          int                   `json:"id"`
	Title       string                `json:"title"`
	Categories  []repository.Category `json:"categories"`
	ContentHash string                `json:"contentHash"`
	ImageHash   string                `json:"imageHash"`
}

type DraftRequestBody struct {
	Article  DraftRequestArticle `json:"article"`
	Contents string              `json:"contents"`
}

func DraftHandler(c *gin.Context) {
	var (
		body RequestBody
		err  error
	)

	if err = ParseRequestBody(c.Request, &body); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fp := ResolveContentFilePath(body.Article.ContentHash, "drafts")
	draft := repository.Draft{
		Id:          body.Article.Id,
		Title:       body.Article.Title,
		Categories:  resolveToDraftCategories(body.Article.Categories),
		UpdateDate:  time.Now(),
		ContentHash: ConvertPathToFileName(fp),
		ImageHash:   body.Article.ImageHash,
	}

	if err = OutputFile(fp, body.Contents); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mysql := repository.GlobalMysql
	if err = repository.DraftCmd(mysql, draft); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		DeleteFile(fp)
		return
	}

	fmt.Fprint(c.Writer, http.StatusOK)
}

func resolveToDraftCategories(categories []repository.Category) string {
	res := ""
	for i, c := range categories {
		if i != 0 {
			res += "&"
		}
		res += c.Name
	}
	return res
}
