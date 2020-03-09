package main

import (
	"flag"
	"net/http"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/repository"
	"github.com/gin-gonic/gin"
)

var (
	Logger  = argus.Logger
	Errors  = &argus.Errors
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
	router.Use(HandleError())

	find := router.Group("/find")
	{
		find.GET("/article/list", handler.FindArticleController)
		find.GET("/article/list/title", handler.FindArticleByTitleController)
		find.GET("/article/list/create-date", handler.FindArticleByCreateDateController)
		find.GET("/article/list/category", handler.FindArticleByCategoryController)
		find.GET("/category/list", handler.FindCategoryController)
		find.GET("/draft/list", handler.FindDraftController)
	}

	register := router.Group("/register")
	{
		register.POST("/article", handler.RegisterArticleController)
	}

	update := router.Group("/update", handler.UpdateArticleController)
	{
		update.PUT("/article")
	}

	save := router.Group("/draft", handler.DraftController)
	{
		save.POST("/article")
	}

	return router
}

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(*Errors) != 0 {
			Logger.ErrorLog(*Errors)
			*Errors = []argus.Error{}
			(c.Writer).WriteHeader(http.StatusInternalServerError)
		}
	}
}
