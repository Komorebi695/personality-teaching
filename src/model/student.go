package model

//学生列表

type StudentListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" binding:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" binding:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" binding:"required"` //每页条数
}

type StudentListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数"  binding:""` //总数
	List  []StudentListItemOutput `json:"list" form:"list" comment:"列表"  binding:""`   //列表
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
	StudentId string `json:"student_id" form:"student_id" comment:"学生ID" binding:"required"` //服务ID
}

//学生添加

type StudentAddInput struct {
	StudentId   string `json:"student_id" form:"student_id" binding:"required"`     //学号
	StudentName string `json:"student_name" form:"student_name" binding:"required"` //学生姓名
	Password    string `json:"password" form:"password" binding:"required"`         //学生密码
	College     string `json:"college" form:"college" binding:"required"`           //学院名称
	Major       string `json:"major" form:"major" binding:"required"`               //专业
	ClassId     string `json:"class_id" form:"class_id" binding:"required"`         //班级编号
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required"` //电话
}

//学生修改

type StudentUpdateInput struct {
	StudentId   string `json:"student_id" form:"student_id" binding:"required"` //id
	StudentName string `json:"student_name" form:"student_name"`                //学生姓名
	PhoneNumber string `json:"phone_number" form:"phone_number"`                //电话
}
