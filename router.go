package main

import "github.com/gin-gonic/gin"

func newRouter() *gin.RouterGroup {
	router := gin.Default().Group("/argus")

	find := router.Group("/find")
	{
		find.GET("/article/list")
		find.GET("/article/category/:category")
		find.GET("/article/date/:date")
		find.GET("/article/title/:title")
		find.GET("/category/list")
	}

	register := router.Group("/register")
	{
		register.POST("/article")
	}

	update := router.Group("/update")
	{
		update.PUT("/article/id/:id")
	}

	return router
}
