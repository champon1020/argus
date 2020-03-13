package handler

import (
	"net/http"
	"path/filepath"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/service"
	"github.com/gin-gonic/gin"
)

func DeleteImageController(c *gin.Context) {
	DeleteImageHandler(c)
}

func DeleteImageHandler(c *gin.Context) {
	imgName := c.Query("imgName")
	fp := filepath.Join(argus.EnvVars.Get("resource"), imgName)
	if err := service.DeleteFile(fp); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
