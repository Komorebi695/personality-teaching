package controller

import (
	"github.com/gin-gonic/gin"
	"personality-teaching/src/logic"
	"personality-teaching/src/middleware"
	"personality-teaching/src/model"
)

type StudentController struct{}

func StudentRegister(group *gin.RouterGroup) {
	student := &StudentController{}
	group.GET("/student_list", student.StudentList)
	group.GET("/student_delete", student.StudentDelete)
	group.POST("/student_add", student.StudentAdd)
}

var studentService logic.StudentService

// StudentList godoc
// @Summary 学生列表
// @Description 学生列表
// @Tags 学生管理
// @ID /student/student_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=model.StudentListOutput} "success"
// @Router /student/student_list [get]
func (student *StudentController) StudentList(c *gin.Context) {
	//从上下文获取参数并校验
	params := &model.StudentListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3000, err)
		return
	}
	out, err := studentService.StudentList(c, params)
	if err != nil {
		middleware.ResponseError(c, 3001, err)
	}
	middleware.ResponseSuccess(c, out)
}

// StudentDelete godoc
// @Summary 学生删除
// @Description 学生删除
// @Tags 学生管理
// @ID /student/student_delete
// @Accept  json
// @Produce  json
// @Param student_id query string true "学生编号"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /student/student_delete [get]
func (student *StudentController) StudentDelete(c *gin.Context) {
	params := &model.StudentDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3000, err)
		return
	}
	err := studentService.StudentDelete(c, params)
	if err != nil {
		middleware.ResponseError(c, 3001, err)
	}
	middleware.ResponseSuccess(c, "")
}

// StudentAdd godoc
// @Summary 添加学生
// @Description 添加学生
// @Tags 学生管理
// @ID /student/student_add
// @Accept  json
// @Produce  json
// @Param body body model.StudentAddInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /student/student_add [post]
func (student *StudentController) StudentAdd(c *gin.Context) {
	params := &model.StudentAddInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3000, err)
		return
	}
	err := studentService.StudentAdd(c, params)
	if err != nil {
		middleware.ResponseError(c, 3001, err)
	}
	middleware.ResponseSuccess(c, "")
}
