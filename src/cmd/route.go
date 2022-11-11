package main

import (
	"personality-teaching/src/controller"
	"personality-teaching/src/middle"

	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	// 这个可以这样-> "/项目前缀/teacher/login"
	router.POST("/teacher/login", controller.TeacherLogin)

	//项目前缀可以加在teacher前面，即 -> router.Group("/项目前缀/teacher")
	//开启登录认证,以下接口需要认证成功才能访问
	teacherRouter := router.Group("/teacher")
	teacherRouter.Use(middle.VerifyTeacher)
	{
		// 班级管理
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
		//题目管理
		teacherRouter.GET("/question/list", controller.QuestionList)
		teacherRouter.DELETE("/question", controller.QuestionDelete)
		teacherRouter.GET("/question/detail", controller.QuestionDetail)
		teacherRouter.POST("/question", controller.QuestionAdd)
		teacherRouter.PUT("/question", controller.QuestionUpdate)
		//知识点管理
		teacherRouter.GET("/point/list", controller.PointList)
		teacherRouter.DELETE("/point", controller.PointDelete)
		teacherRouter.GET("/point/detail", controller.PointDetail)
		teacherRouter.POST("/point", controller.PointAdd)
		teacherRouter.PUT("/point", controller.PointUpdate)
	}

	return router
}
