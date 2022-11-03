package main

import (
	"personality-teaching/src/controller"
	"personality-teaching/src/middle"

	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	preRouter := router.Group("") //项目前缀，先空着

	preRouter.POST("/teacher/login", controller.TeacherLogin)

	//  班级管理
	preRouter.POST("/teacher/class", middle.VerifyTeacher, controller.AddClass)
	preRouter.PUT("/teacher/class", middle.VerifyTeacher, controller.UpdateClass)
	preRouter.DELETE("/teacher/class", middle.VerifyTeacher, controller.DeleteClass)
	preRouter.GET("/teacher/class", controller.ClassInfo)
	preRouter.GET("/teacher/class/list", middle.VerifyTeacher, controller.ClassList)

	//题目模块接口
	questionGroup := router.Group("/question")
	questionGroup.Use(middlewares...)
	{
		controller.QuestionRegister(questionGroup)
	}
	//题目分类模块接口
	QuestionTypeGroup := router.Group("/type")
	QuestionTypeGroup.Use(middlewares...)
	{
		controller.QuestionTypeRegister(QuestionTypeGroup)
	}
	//知识点模块接口
	knowledgePointGroup := router.Group("/point")
	knowledgePointGroup.Use(middlewares...)
	{
		controller.KnowledgePointRegister(knowledgePointGroup)
	}
	return router
}
