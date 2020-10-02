package main

import (
	"fmt"
	"os"
	"time"

	"github.com/champon1020/argus/v2"
	"github.com/champon1020/argus/v2/handler"
	"github.com/champon1020/argus/v2/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	argus.Init()
	model.InitDatabase()
}

func main() {
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

	find := r.Group("/api/find")
	{
		find.GET("/article/list", wrapHandler(handler.FindArticlesList))
	}
	return r
}

func wrapHandler(h func(c *gin.Context, db model.DatabaseIface) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c, model.Db)
		if err != nil {
			// parse error
		}
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
		},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}
}
