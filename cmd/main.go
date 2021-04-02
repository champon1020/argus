package main

import (
	"log"
	"os"

	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/interfaces/di"
	"github.com/champon1020/argus/interfaces/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	di := di.NewDI(config)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://blog.champonian.com", "http://localhost:8000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[ECHO] ${method} | ${uri} | ${status} : ${error}\n",
		Output: os.Stdout,
	}))

	router.AppRouter(e, di.NewAppHandler())
	e.Logger.Fatal(e.Start(":8000"))
}
