package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kastnerorz/animal-crossing-trading-system/backend/middlewares"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CorsMiddleware())
	authMiddleware := middlewares.AuthMiddleware()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", CreateUser)
		v1.GET("/users/:id", GetUser)
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.POST("/refresh", authMiddleware.RefreshHandler)
		v1.GET("/quotations", GetQuotations)
		v1.GET("/quotations/:id", GetQuotation)
		v1.Use(authMiddleware.MiddlewareFunc())
		{
			v1.GET("/me", GetMyInfo)
			v1.PUT("/me", UpdateUser)
			v1.POST("/quotations", CreateQuotation)
			v1.PUT("/quotations/:id", UpdateQuotation)
			v1.GET("/my-quotations", GetMyQuotation)
			v1.DELETE("/quotations/:id", DeleteQuotation)
			v1.POST("/applications", CreateApplication)
			v1.GET("/applications", GetMyApplications)
			v1.PUT("/applications/:id", UpdateApplication)
			v1.DELETE("/applications/:id", DeleteApplication)
		}

	}

	router.Use(gin.Recovery())
	return router
}
