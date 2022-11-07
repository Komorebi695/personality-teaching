package main

import (
	"personality-teaching/src/controller"
	"personality-teaching/src/middle"

	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	// 这个可以这样-> "/项目前缀/login"
	router.POST("/login", controller.TeacherLogin)

	//项目前缀可以加在teacher前面，即 -> router.Group("/项目前缀/teacher")
	teacherRouter := router.Group("/teacher")
	teacherRouter.Use(middle.VerifyTeacher)
	{
		//  班级管理
		teacherRouter.POST("/class", controller.AddClass)
		teacherRouter.PUT("/class", controller.UpdateClass)
		teacherRouter.DELETE("/class", controller.DeleteClass)
		teacherRouter.GET("/class", controller.ClassInfo)
		teacherRouter.GET("/class/list", controller.ClassList)
		// 试卷管理
		teacherRouter.POST("/exam", controller.AddExam)
		teacherRouter.PUT("/exam", controller.UpdateExam)
		teacherRouter.DELETE("/exam", controller.DeleteExam)
		teacherRouter.GET("/exam", controller.ExamInfo)
		teacherRouter.GET("/exam/list", controller.ExamList)
		teacherRouter.POST("/exam/send/:id", controller.SendExam)
	}

	return router
}
