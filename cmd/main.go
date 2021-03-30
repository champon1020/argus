package main

import (
	"fmt"
	"os"

	"github.com/champon1020/argus/interfaces/di"
	"github.com/champon1020/argus/interfaces/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func dns() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("ARGUS_DB_USER"),
		os.Getenv("ARGUS_DB_PASSWORD"),
		os.Getenv("ARGUS_DB_HOST"),
		os.Getenv("ARGUS_DB_PORT"),
		os.Getenv("ARGUS_DB_NAME"),
	)
}

func main() {
	di := di.NewDI()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://blog.champonian.com", "http://localhost:8000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[ECHO] ${method} | ${uri} | ${status}\n",
		Output: os.Stdout,
	}))

	router.AppRouter(e, di.NewAppHandler())
	e.Logger.Fatal(e.Start(":8000"))
}
