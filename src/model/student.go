package model

import (
	"github.com/gin-gonic/gin"
	"personality-teaching/src/public"
)

//学生列表

type StudentListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` //每页条数
}

//输入参数校验

func (param *StudentListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type StudentListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数"  validate:""` //总数
	List  []StudentListItemOutput `json:"list" form:"list" comment:"列表"  validate:""`   //列表
}

type StudentListItemOutput struct {
	StudentID   string `json:"student_id" form:"student_id"`     //id
	StudentName string `json:"student_name" form:"student_name"` //学生姓名
	College     string `json:"college" form:"college"`           //学院名称
	Major       string `json:"major" form:"major"`               //专业
	Class       string `json:"class" form:"class"`               //班级
	PhoneNumber string `json:"phone_number" form:"phone_number"` //电话
}

//学生删除

type StudentDeleteInput struct {
	StudentId string `json:"student_id" form:"student_id" comment:"学生ID" validate:"required"` //服务ID
}

func (param *StudentDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

//学生添加

type StudentAddInput struct {
	StudentId   string `json:"student_id" form:"student_id" validate:"required"`     //学号
	StudentName string `json:"student_name" form:"student_name" validate:"required"` //学生姓名
	Password    string `json:"password" form:"password" validate:"required"`         //学生密码
	College     string `json:"college" form:"college" validate:"required"`           //学院名称
	Major       string `json:"major" form:"major" validate:"required"`               //专业
	ClassId     string `json:"class_id" form:"class_id" validate:"required"`         //班级编号
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"` //电话
}

func (param *StudentAddInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

//学生修改

type StudentUpdateInput struct {
	StudentId   string `json:"student_id" form:"student_id" validate:"required"` //id
	StudentName string `json:"student_name" form:"student_name"`                 //学生姓名
	PhoneNumber string `json:"phone_number" form:"phone_number"`                 //电话
}

func (param *StudentUpdateInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
