package main

import (
	"os"

	"github.com/champon1020/argus"
	"github.com/champon1020/argus/model"
	"github.com/champon1020/argus/route"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	argus.Init()
	model.InitDatabase()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://blog.champonian.com", "http://localhost:8000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[ECHO] ${method} | ${uri} | ${status}\n",
		Output: os.Stdout,
	}))

	route.AddRoutes(e)
	e.Logger.Fatal(e.Start(":8000"))
}
