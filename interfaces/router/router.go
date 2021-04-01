package router

import (
	"errors"

	"github.com/champon1020/argus/interfaces/handler"
	"github.com/labstack/echo/v4"
)

func dummyHandler(c echo.Context) error {
	return errors.New("Not implemented")
}

// AppRouter declares the api routes.
func AppRouter(e *echo.Echo, h handler.AppHandler) {
	v3 := e.Group("/api/v3")
	{
		v3.GET("/get/article/id/:id", h.PublicArticleByID)
		v3.GET("/get/articles", h.PublicArticles)
		v3.GET("/get/articles/title/:title", h.PublicArticlesByTitle)
		v3.GET("/get/articles/tag/:tag", h.PublicArticlesByTag)
		v3.GET("/get/tags", h.PublicTags)
	}

	private := v3.Group("/private")
	{
		private.GET("/get/article/id/:id", h.ArticleByID)
		private.GET("/get/articles", h.Articles)
		private.GET("/get/images", dummyHandler)
		private.POST("/post/article", h.PostArticle)
		private.POST("/post/image", dummyHandler)
		private.PUT("/update/article", h.UpdateArticle)
		private.PUT("/update/article/isPublic", dummyHandler)
		private.DELETE("/delete/article", h.DeleteArticle)
		private.DELETE("/delete/image", dummyHandler)

		private.POST("/verify/token", dummyHandler)
	}
}
