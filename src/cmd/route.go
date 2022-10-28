package main

import (
	"personality-teaching/src/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	preRouter := router.Group("") //项目前缀，先空着

	preRouter.POST("/teacher/login", controller.TeacherLogin)
	return router
}
