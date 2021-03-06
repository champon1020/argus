package router

import (
	"net/http"

	"github.com/champon1020/argus/interfaces/handler"
	"github.com/champon1020/argus/interfaces/middleware"
	"github.com/labstack/echo/v4"
)

// AppRouter declares the api routes.
func AppRouter(e *echo.Echo, h handler.AppHandler) {
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	v3 := e.Group("/api/v3")
	{
		v3.GET("/get/article/id/:id", h.PublicArticleByID)
		v3.GET("/get/articles", h.PublicArticles)
		v3.GET("/get/articles/title/:title", h.PublicArticlesByTitle)
		v3.GET("/get/articles/tag/:tag", h.PublicArticlesByTag)
		v3.GET("/get/tags", h.PublicTags)
		v3.GET("/get/headerImages", h.HeaderImages)
	}

	private := v3.Group("/private")
	private.Use(middleware.AuthMiddleware)
	{
		private.GET("/get/article/id/:id", h.ArticleByID)
		private.GET("/get/articles", h.Articles)
		private.GET("/get/drafts", h.DraftArticles)
		private.GET("/get/images", h.Images)
		private.POST("/post/article", h.PostArticle)
		private.POST("/post/image", h.PostImage)
		private.PUT("/update/article", h.UpdateArticle)
		private.PUT("/update/article/status", h.UpdateArticleStatus)
		private.DELETE("/delete/article", h.DeleteArticle)
		private.DELETE("/delete/images", h.DeleteImage)
		private.POST("/verify", h.VerifyToken)
	}
}
