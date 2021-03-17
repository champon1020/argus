package main

import (
	"fmt"
	"log"
	"os"

	"github.com/champon1020/argus/handler"
	"github.com/champon1020/argus/repository"
	"github.com/champon1020/argus/route"
	"github.com/champon1020/mgorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
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
	db, err := mgorm.New("mysql", dns())
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	h := handler.NewHandler(repo)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://blog.champonian.com", "http://localhost:8000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[ECHO] ${method} | ${uri} | ${status}\n",
		Output: os.Stdout,
	}))

	route.AddRoutes(e, h)
	e.Logger.Fatal(e.Start(":8000"))
}
