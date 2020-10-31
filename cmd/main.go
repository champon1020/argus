package main

import (
	"fmt"
	"os"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/argus-private/auth"
	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/handler/private"
	"github.com/champon1020/argus/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	argus.Init()
	model.InitDatabase()

	r := newRouter()
	_ = r.Run(":8000")
}

func newRouter() *gin.Engine {
	r := gin.New()

	// Set the loggin configuration
	r.Use(gin.LoggerWithConfig(*loggerConfig()))

	r.Use(gin.Recovery())

	// Set the cors configuration
	r.Use(cors.New(*corsConfig()))

	r.GET("/healthcheck", func(c *gin.Context) {
		c.AbortWithStatus(200)
		return
	})

	find := r.Group("/api/find")
	{
		find.GET("/article/list", wrapHandlerWithDatabase(handler.APIFindArticles))
		find.GET("/article/id", wrapHandlerWithDatabase(handler.APIFindArticlesByID))
		find.GET("/article/list/title", wrapHandlerWithDatabase(handler.APIFindArticlesByTitle))
		find.GET("/article/list/category", wrapHandlerWithDatabase(handler.APIFindArticlesByCategory))
		find.GET("/category/list", wrapHandlerWithDatabase(handler.APIFindCategories))
	}

	r.POST("/api/verify/token", wrapHandler(auth.APIVerify))

	priv := r.Group("/api/private")
	priv.Use(wrapHandler(auth.Middleware))
	{
		find := priv.Group("/find")
		{
			find.GET("/article/id", wrapHandlerWithDatabase(private.APIFindArticleByID))
			find.GET("/article/list", wrapHandlerWithDatabase(private.APIFindArticles))
			find.GET("/draft/id", wrapHandlerWithDatabase(private.APIFindDraftByID))
			find.GET("/draft/list", wrapHandlerWithDatabase(private.APIFindDrafts))
			find.GET("/image/list", wrapHandlerWithDatabase(private.APIFindImages))
		}

		register := priv.Group("/register")
		{
			register.POST("/article", wrapHandlerWithDatabase(private.APIRegisterArticle))
			register.POST("/draft", wrapHandlerWithDatabase(private.APIRegisterDraft))
			register.POST("/image", wrapHandlerWithDatabase(private.APIRegisterImage))
		}

		update := priv.Group("/update")
		{
			update.PUT("/article", wrapHandlerWithDatabase(private.APIUpdateArticle))
			update.PUT("/article/isPrivate", wrapHandlerWithDatabase(private.APIUpdateIsPrivate))
			update.PUT("/draft", wrapHandlerWithDatabase(private.APIUpdateDraft))
		}

		delete := priv.Group("/delete")
		{
			delete.DELETE("/draft", wrapHandlerWithDatabase(private.APIDeleteDraft))
			delete.DELETE("/image", wrapHandlerWithDatabase(private.APIDeleteImage))
		}
	}

	return r
}

func wrapHandler(h func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c)

		if err != nil {
			printErr(err)
		}
	}
}

func wrapHandlerWithDatabase(h func(c *gin.Context, db model.DatabaseIface) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c, model.Db)
		if err != nil {
			printErr(err)
		}
	}
}

// printErr outputs the error log as standard output.
func printErr(err error) {
	if e, ok := err.(*argus.Error); ok {
		e.Log()
	}
}

func loggerConfig() *gin.LoggerConfig {
	return &gin.LoggerConfig{
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
	}
}

func corsConfig() *cors.Config {
	return &cors.Config{
		AllowAllOrigins: false,
		AllowOrigins: []string{
			"https://blog.champonian.com",
			//"http://localhost:3000",
		},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}
}
