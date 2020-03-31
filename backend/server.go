package main

import (
	"flag"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.StringVar(&MongoURI, "mongo-url", "mongodb://localhost:27017", "")
	flag.StringVar(&MongoCollection, "mongo-collection", "acts-dev", "")
	flag.StringVar(&Port, "port", "8080", "")
	flag.Parse()
	authMiddleware := AuthMiddleware()

	router := gin.Default()
	router.Use(CorsMiddleware())

	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", CreateUser)
		v1.GET("/users/:id", GetUser)
		//v1.PUT("/users/:id", UpdateUser)
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.POST("/refresh", authMiddleware.RefreshHandler)
		v1.GET("/quotations", GetQuotations)
		v1.Use(authMiddleware.MiddlewareFunc())
		{
			v1.GET("/me", GetMyInfo)
			v1.POST("/quotations", CreateQuotation)
			v1.PUT("/quotations/:id", UpdateQuotation)
			v1.GET("/quotations/my", GetMyQuotation)
			v1.POST("/applications", CreateApplication)
			v1.GET("/applications", GetMyApplications)
			v1.PUT("/applications/:id", UpdateApplication)
		}

	}

	router.Use(gin.Recovery())
	router.Run(":" + Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
