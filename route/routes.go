package route

import "github.com/labstack/echo"

func dummyHandler(c echo.Context) error {
	return nil
}

// AddRoutes registers new routes.
func AddRoutes(e *echo.Echo) {
	v3 := e.Group("/api/v3")

	{
		v3.GET("/get/article/id/:id", dummyHandler)
		v3.GET("/get/articles", dummyHandler)
		v3.GET("/get/categories", dummyHandler)
		v3.GET("/get/articles", dummyHandler)
	}

	private := v3.Group("/private")
	{
		private.GET("/get/article/id/:id", dummyHandler)
		private.GET("/get/draft/id/:id", dummyHandler)
		private.GET("/get/articles", dummyHandler)
		private.GET("/get/drafts", dummyHandler)
		private.GET("/get/images", dummyHandler)
		private.GET("/count/articles", dummyHandler)
		private.GET("/count/drafts", dummyHandler)
		private.GET("/count/images", dummyHandler)
		private.POST("/post/article", dummyHandler)
		private.POST("/post/draft", dummyHandler)
		private.POST("/post/image", dummyHandler)
		private.POST("/verify/token", dummyHandler)
		private.PUT("/update/article", dummyHandler)
		private.PUT("/update/article/isPublic", dummyHandler)
		private.PUT("/update/draft", dummyHandler)
		private.DELETE("/delete/article", dummyHandler)
		private.DELETE("/delete/draft", dummyHandler)
		private.DELETE("/delete/image", dummyHandler)
	}
}
