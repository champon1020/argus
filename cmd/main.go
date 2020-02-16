package main

import (
	"github.com/champon1020/argus"
	"github.com/champon1020/argus/handler"
	"github.com/gin-gonic/gin"
)

var logger argus.Logger

func init() {
	logger.NewLogger("[main]")
}

func main() {
	r := NewRouter()
	r.Run(":8080")
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	find := router.Group("/find")
	{
		find.GET("/article/list", handler.FindArticleHandler)
		find.GET("/article/category/:category")
		find.GET("/article/date/:date")
		find.GET("/article/title/:title")
		find.GET("/category/list", handler.FindCategoryHandler)
	}

	register := router.Group("/register")
	{
		register.POST("/article", handler.RegisterArticleHandler)
	}

	update := router.Group("/update", handler.UpdateArticleHandler)
	{
		update.PUT("/article/id/:id")
	}

	return router
}
