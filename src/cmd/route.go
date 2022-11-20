package main

import (
	"github.com/gin-gonic/gin"
	"personality-teaching/src/controller"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	// 这个可以这样-> "/项目前缀/teacher/login"
	router.POST("/teacher/login", controller.TeacherLogin)

	//项目前缀可以加在teacher前面，即 -> router.Group("/项目前缀/teacher")
	//开启登录认证,以下接口需要认证成功才能访问
	teacherRouter := router.Group("/teacher")
	//teacherRouter.Use(middle.VerifyTeacher)
	{
		// 班级管理
		teacherRouter.POST("/class", controller.AddClass)
		teacherRouter.PUT("/class", controller.UpdateClass)
		teacherRouter.DELETE("/class", controller.DeleteClass)
		teacherRouter.GET("/class", controller.ClassInfo)
		teacherRouter.GET("/class/list", controller.ClassList)
		// 班级学生管理
		teacherRouter.POST("/student", controller.CreateStudent)
		teacherRouter.POST("/class/student", controller.AddStudentToClass)
		teacherRouter.GET("/class/student/list", controller.StudentsInClass)
		teacherRouter.GET("/student/list", controller.StudentNotInClass)
		teacherRouter.DELETE("/class/student", controller.DeleteClassStudent)
		// 试卷管理
		teacherRouter.POST("/exam", controller.AddExam)
		teacherRouter.PUT("/exam", controller.UpdateExam)
		teacherRouter.DELETE("/exam", controller.DeleteExam)
		teacherRouter.GET("/exam", controller.ExamInfo)
		teacherRouter.GET("/exam/list", controller.ExamList)
		teacherRouter.POST("/exam/send/:id", controller.SendExam)
		teacherRouter.POST("/exam/search", controller.SearchExam)
		//试卷评阅
		teacherRouter.GET("/review/class", controller.ReviewClass)
		teacherRouter.GET("/review/student/list", controller.ReviewStudentList)
		teacherRouter.GET("/review/student", controller.ReviewStudent)
		teacherRouter.PUT("/review", controller.ReviewUpdate)

		//题目管理
		teacherRouter.GET("/question/list", controller.QuestionList)
		teacherRouter.DELETE("/question", controller.QuestionDelete)
		teacherRouter.GET("/question/detail", controller.QuestionDetail)
		teacherRouter.POST("/question", controller.QuestionAdd)
		teacherRouter.PUT("/question", controller.QuestionUpdate)
		//知识点管理
		teacherRouter.GET("/point/list", controller.PointList)
		teacherRouter.GET("/point/list/one_stage", controller.PointOneStageList)
		teacherRouter.DELETE("/point", controller.PointDelete)
		teacherRouter.GET("/point/detail", controller.PointDetail)
		teacherRouter.POST("/point", controller.PointAdd)
		teacherRouter.PUT("/point", controller.PointUpdate)
	}

	return router
}
