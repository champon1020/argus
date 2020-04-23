package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/champon1020/argus/argus-private/auth"

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
	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[GIN] %s | %s |%s %d %s| %15s | %15s |%s %s %s %s \n",
				param.Request.Proto,
				param.TimeStamp.Format("2006-01-02 15:04:05 MST -0700"),
				param.StatusCodeColor(),
				param.StatusCode,
				param.ResetColor(),
				param.Latency,
				param.ClientIP,
				param.MethodColor(),
				param.Method,
				param.ResetColor(),
				param.Path,
			)
		},
		Output:    os.Stdout,
		SkipPaths: []string{"/healthcheck"},
	}))

	router.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowAllOrigins: false,
		AllowOrigins: []string{
			"https://blog.champonian.com",
		},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}

	if argus.EnvVars.Get("mode") == "dev" {
		corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, "http://localhost:3000")
	}

	router.Use(cors.New(corsConfig))
	router.Use(HandleError())

	router.GET("/healthcheck", HealthCheck)

	find := router.Group("/api/find")
	{
		find.GET("/article/pickup", handler.FindPickUpArticleController)
		find.GET("/article/sortedId", handler.FindArticleBySortedIdController)
		find.GET("/article/list", handler.FindArticleController)
		find.GET("/article/list/title", handler.FindArticleByTitleController)
		find.GET("/article/list/create-date", handler.FindArticleByCreateDateController)
		find.GET("/article/list/category", handler.FindArticleByCategoryController)
		find.GET("/category/list", handler.FindCategoryController)
	}

	router.POST("/api/verify/token", auth.VerifyHandler)

	private := router.Group("/api/private")
	private.Use(auth.Middleware)
	{
		find := private.Group("/find")
		{
			find.GET("/article/id", handler.FindArticleByIdController)
			find.GET("/draft/id", handler.FindDraftByIdController)
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
		}

		delete := private.Group("/delete")
		{
			delete.DELETE("/draft", handler.DeleteDraftController)
			delete.DELETE("/image", handler.DeleteImageController)
		}

		draft := private.Group("/draft")
		{
			draft.POST("/article", handler.DraftController)
		}
	}

	return router
}

func HealthCheck(c *gin.Context) {
	c.AbortWithStatus(200)
	return
}

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(*Errors) > 0 {
			Logger.ErrorLog(*Errors)
			*Errors = []argus.Error{}
			(c.Writer).WriteHeader(http.StatusInternalServerError)
		}
	}
}
