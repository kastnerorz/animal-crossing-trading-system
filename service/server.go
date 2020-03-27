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
			v1.POST("/quotations", CreateQuotation)
			v1.PUT("/quotations/:id", UpdateQuotation)
		}

	}

	//router.Use(gin.Recovery())
	router.Run(":" + Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
