package main

import (
	"log"

	"github.com/champon1020/argus"
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

	logger := argus.NewLogger()

	di := di.NewDI(config, logger)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://champonian.com", "http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	router.AppRouter(e, di.NewAppHandler())
	e.Logger.Fatal(e.Start(":8000"))
}
