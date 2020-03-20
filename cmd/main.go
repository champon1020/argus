package main

import (
	"net/http"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/repo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	Logger = argus.Logger
	Errors = &argus.Errors
	r      *gin.Engine
)

func init() {
	repo.GlobalMysql = repo.NewMysql()
	r = NewRouter()
}

func main() {
	_ = r.Run(":8000")
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"http://localhost:3000"},
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))
	router.Use(HandleError())

	find := router.Group("/api/find")
	{
		find.GET("/article/pickup", handler.FindPickUpArticleController)
		find.GET("/article/id", handler.FindArticleByIdController)
		find.GET("/article/sortedId", handler.FindArticleBySortedIdController)
		find.GET("/article/list", handler.FindArticleController)
		find.GET("/article/list/title", handler.FindArticleByTitleController)
		find.GET("/article/list/create-date", handler.FindArticleByCreateDateController)
		find.GET("/article/list/category", handler.FindArticleByCategoryController)
		find.GET("/category/list", handler.FindCategoryController)
	}

	private := router.Group("/api/private")
	{
		find := private.Group("/find")
		{
			find.GET("/article/list/all", handler.FindAllArticleController)
			find.GET("/draft/list", handler.FindDraftController)
			find.GET("/image/list", handler.FindImageController)
		}

		register := private.Group("/register")
		{
			register.POST("/article", handler.RegisterArticleController)
			register.POST("/image", handler.RegisterImageController)
		}

		update := private.Group("/update")
		{
			update.PUT("/article", handler.UpdateArticleController)
			update.PUT("/article/object", handler.UpdateArticleObjController)
		}

		delete := private.Group("/delete")
		{
			delete.DELETE("/image", handler.DeleteImageController)
		}

		draft := private.Group("/draft")
		{
			draft.POST("/article", handler.DraftController)
		}
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
