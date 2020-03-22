package handler

import (
	"net/http"
	"path/filepath"

	"github.com/champon1020/argus/repo"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

func DeleteImageController(c *gin.Context) {
	_ = DeleteImageHandler(c)
}

func DeleteImageHandler(c *gin.Context) (err error) {
	imgNames := c.QueryArray("imageNames")

	for _, name := range imgNames {
		fp := filepath.Join(argus.EnvVars.Get("resource"), "images", name)
		if err = service.DeleteFile(fp); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	c.Writer.WriteHeader(http.StatusOK)
	return
}

func DeleteDraftController(c *gin.Context) {
	_ = DeleteDraftHandler(c, repo.DeleteDraftCommand)
}

func DeleteDraftHandler(c *gin.Context, repoCmd repo.DeleteDraftCmd) (err error) {
	id := c.Query("id")
	hash := c.Query("contentHash")
	draft := repo.Draft{Id: id}

	fp := filepath.Join(argus.EnvVars.Get("resource"), "drafts", hash+"_md")
	if err = service.DeleteFile(fp); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = repoCmd(*repo.GlobalMysql, draft); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	return
}
