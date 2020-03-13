package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	Logger = argus.Logger
	Errors = &argus.Errors
)

func main() {
	flag.Parse()
	argus.GlobalConfig = argus.NewConfig(flag.Arg(0))
	repository.GlobalMysql = repository.NewMysql()

	r := NewRouter()
	_ = r.Run(":8080")
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"http://localhost:3000"},
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))
	router.Use(HandleError())

	find := router.Group("/api/find")
	{
		find.GET("/article/list", handler.FindArticleController)
		find.GET("/article/list/id", handler.FindArticleByIdController)
		find.GET("/article/list/title", handler.FindArticleByTitleController)
		find.GET("/article/list/create-date", handler.FindArticleByCreateDateController)
		find.GET("/article/list/category", handler.FindArticleByCategoryController)
		find.GET("/category/list", handler.FindCategoryController)
		find.GET("/draft/list", handler.FindDraftController)
		find.GET("/image/list")
	}

	register := router.Group("/api/register")
	{
		register.POST("/article", handler.RegisterArticleController)
		register.POST("/image", handler.RegisterImageController)
	}

	update := router.Group("/api/update")
	{
		update.PUT("/article", handler.UpdateArticleController)
	}

	delete := router.Group("/api/delete")
	{
		delete.DELETE("/image")
	}

	draft := router.Group("/api/draft")
	{
		draft.POST("/article", handler.DraftController)
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
