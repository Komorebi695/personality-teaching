package main

import (
	"github.com/gin-gonic/gin"
	"personality-teaching/src/Cos"
	"personality-teaching/src/controller"
	"personality-teaching/src/middle"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	// 这个可以这样-> "/项目前缀/teacher/login"
	router.POST("/teacher/login", controller.TeacherLogin)
	router.PUT("/pwd", middle.VerifyAny, middle.ChangePassword)
	// 项目前缀可以加在teacher前面，即 -> router.Group("/项目前缀/teacher")
	// 开启登录认证,以下接口需要认证成功才能访问
	teacherRouter := router.Group("/teacher")
	teacherRouter.Use(middle.VerifyTeacher)
	{
		//教师信息
		teacherRouter.GET("/info", controller.TeacherInfo)
		// 班级管理
		teacherRouter.POST("/class", controller.AddClass)
		teacherRouter.PUT("/class", controller.UpdateClass)
		teacherRouter.DELETE("/class", controller.DeleteClass)
		teacherRouter.GET("/class", controller.ClassInfo)
		teacherRouter.GET("/class/list", controller.ClassList)
		teacherRouter.GET("/class/check", controller.ClassNameCheck)
		// 班级学生管理
		teacherRouter.POST("/class/student", controller.AddStudentToClass)
		teacherRouter.GET("/class/student/list", controller.StudentsInClass)
		teacherRouter.GET("/student/list", controller.StudentNotInClass)
		teacherRouter.DELETE("/class/student", controller.DeleteClassStudent)
		// 学生管理
		teacherRouter.POST("/student", controller.CreateStudent)
		teacherRouter.GET("/student/search", controller.SearchStudent)
		teacherRouter.DELETE("/student", controller.DeleteStudent)
		teacherRouter.PUT("/student", controller.UpdateStudent)

		// 试卷管理
		teacherRouter.POST("/exam", controller.AddExam)
		teacherRouter.PUT("/exam", controller.UpdateExam)
		teacherRouter.DELETE("/exam", controller.DeleteExam)
		teacherRouter.GET("/exam", controller.ExamInfo)
		teacherRouter.GET("/exam/list", controller.ExamList)
		teacherRouter.POST("/exam/send/:id", controller.SendExam)
		teacherRouter.POST("/exam/search", controller.SearchExam)
		teacherRouter.GET("/exam/student/list", controller.ReleaseStudentList)
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
		teacherRouter.POST("Cquestion", Cos.QuestionUploadFileToCos)
		//知识点管理
		teacherRouter.GET("/point/list", controller.PointList)
		teacherRouter.GET("/point/list/one_stage", controller.PointOneStageList)
		teacherRouter.DELETE("/point", controller.PointDelete)
		teacherRouter.GET("/point/detail", controller.PointDetail)
		teacherRouter.POST("/point", controller.PointAdd)
		teacherRouter.PUT("/point", controller.PointUpdate)
		teacherRouter.PUT("/point/connection", controller.PointConnectionUpdate)
		teacherRouter.POST("/point/uploadImage", Cos.KnpUploadFileToCos)
	}

	// 学生登录 123
	studentRouter := router.Group("/student")
	studentRouter.POST("/login", controller.StudentLogin)
	studentRouter.Use(middle.VerifyStudent)
	{ //知识点管理
		studentRouter.GET("/point/list", controller.PointList)
		studentRouter.GET("/point/list/one_stage", controller.PointOneStageList)
		//题目管理
		studentRouter.GET("/question/list", controller.QuestionList)
		studentRouter.GET("/question/detail", controller.QuestionDetail)
		//试卷管理
		studentRouter.GET("/exam/get", controller.GetTeacherExamList)
		studentRouter.PUT("/exam/upload", controller.PostStudentExamAnswer)
		//提交答案回显
		studentRouter.GET("/exam/review", controller.ReviewStudentAnswer)
	}
	return router
}
