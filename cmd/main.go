package main

import (
	"flag"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
)

var (
	Logger  = argus.Logger
	configs argus.Configurations
)

func main() {
	flag.Parse()

	configs.New(flag.Arg(0))
	repository.NewMysql()

	r := NewRouter()
	r.Run(":8080")
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	find := router.Group("/find")
	{
		find.GET("/article/list", handler.FindArticleHandler)
		find.GET("/article/list/title", handler.FindArticleHandlerByTitle)
		find.GET("/article/list/create-date", handler.FindArticleHandlerByCreateDate)
		find.GET("/article/list/category", handler.FindArticleHandlerByCategory)
		find.GET("/category/list", handler.FindCategoryHandler)
	}

	register := router.Group("/register")
	{
		register.POST("/article", handler.RegisterArticleHandler)
	}

	update := router.Group("/update", handler.UpdateArticleHandler)
	{
		update.PUT("/article")
	}

	return router
}
