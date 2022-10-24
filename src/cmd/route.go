package main

import (
	"github.com/gin-gonic/gin"
	"personality-teaching/src/controller"
	"personality-teaching/src/middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//学生模块接口
	studentGroup := router.Group("/student")
	studentGroup.Use(
		middleware.RecoveryMiddleware(),
		middleware.TranslationMiddleware(),
		middleware.RequestLog())
	{
		controller.StudentRegister(studentGroup)
	}
	return router
}
