package model

import "time"

type StudentExams struct {
	ExamID    string  `gorm:"column:exam_id" json:"exam_id"`
	StudentID string  `gorm:"column:student_id" json:"class_id"`
	ExamName  string  `gorm:"column:exam_name" json:"exam_name"`                               // 试卷编号
	Answers   string  `gorm:"column:answers" json:"answers"`                                   // 学生答案
	Score     float64 `gorm:"score" json:"score"`                                              // 分数
	Comment   string  `gorm:"column:comment" form:"comment" binding:"required" json:"comment"` // 备注
}

type ReviewClass struct {
	Count int `gorm:"count" json:"count"`
	BaseClassInfo
}

type ReviewStudent struct {
	Name       string    `gorm:"name" json:"name"`
	Score      float64   `gorm:"score" json:"score"`
	Status     int       `gorm:"status" json:"status"`
	UpdateTime time.Time `gorm:"update_time" json:"update_time"`
}

type ReviewUpdate struct {
	ReviewUpdateReq
	UpdateTime string `gorm:"column:update_time" json:"update_time"` // 更新时间
}

type ReviewUpdateReq struct {
	ExamID    string  `gorm:"column:exam_id" json:"exam_id"`       // 试卷编号
	StudentID string  `gorm:"column:student_id" json:"student_id"` //学生编号
	Answers   string  `gorm:"column:answers" json:"answers"`       // 学生答案
	Score     float64 `gorm:"score" json:"score"`                  // 分数
	Status    int     `gorm:"status" json:"status"`                //状态
}

type ReviewStudentListReq struct {
	ClassID string `json:"class_id"`
	ExamID  string `json:"exam_id"`
}

type ReviewStudentReq struct {
	ExamID    string `json:"exam_id"`
	StudentID string `json:"class_id"`
}
