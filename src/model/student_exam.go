package model

// StudentExam ,个人试卷发放结构
type StudentExam struct {
	SendReq
	StudentID  string `gorm:"column:student_id" json:"student_id"`   // 学生编号
	Answers    string `gorm:"column:answers" json:"answers"`         // 学生答案
	UpdateTime string `gorm:"column:update_time" json:"update_time"` // 更新时间
	CreateTime string `gorm:"column:create_time" json:"create_time"` // 创建时间
}

func (se StudentExam) TableName() string {
	return "t_student_exam"
}

// SendReq ,试卷发放请求基础数据结构
type SendReq struct {
	ExamID    string `gorm:"column:exam_id" form:"exam_id" binding:"required" json:"exam_id"` // 试卷编号
	StartTime string `gorm:"column:start_time" form:"start_time" json:"start_time"`           // 考试开始时间
	EndTime   string `gorm:"column:end_time" form:"end_time" json:"end_time"`                 // 考试结束时间
	Comment   string `gorm:"column:comment" form:"comment" binding:"required" json:"comment"` // 备注
}

// SendPersonReq ,个人
type SendPersonReq struct {
	StudentID string `gorm:"column:student_id" form:"student_id" binding:"required" json:"student_id"` //学生ID
	SendReq
}

// SendClassReq ,班级
type SendClassReq struct {
	ClassID string `gorm:"column:class_id" form:"class_id" binding:"required" json:"class_id"` //班级编号
	SendReq
}

// ClassExam ,班级试卷发放结构
type ClassExam struct {
	SendReq
	ClassID    string `gorm:"column:class_id" form:"class_id" binding:"required" json:"class_id"` //班级编号
	UpdateTime string `gorm:"column:update_time" json:"update_time"`                              // 更新时间
	CreateTime string `gorm:"column:create_time" json:"create_time"`                              // 创建时间
}
