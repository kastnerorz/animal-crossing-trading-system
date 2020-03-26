package main

import (
	"flag"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.StringVar(&MongoURI, "mongo-url", "mongodb://localhost:27017", "")
	flag.StringVar(&MongoCollection, "mongo-collection", "wj-dev", "")
	flag.StringVar(&Port, "port", "8080", "")
	flag.Parse()
	authMiddleware := AuthMiddleware()

	router := gin.Default()
	router.POST("/users", CreateUser)
	router.GET("/users/:id", GetUser)
	//router.PUT("/users/:id", UpdateUser)
	router.Use(authMiddleware.MiddlewareFunc()).POST("/quotations", CreateQuotation)
	router.GET("/quotations", GetQuotations)
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/refresh", authMiddleware.RefreshHandler)

	router.Use(gin.Recovery())
	router.Run(":" + Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
