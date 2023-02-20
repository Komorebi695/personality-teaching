package model

import "time"

type StudentExams struct {
	ExamID        string `gorm:"column:exam_id" json:"exam_id"`
	StudentID     string `gorm:"column:student_id" json:"student_id"`
	ExamName      string `gorm:"column:exam_name" json:"exam_name"`                               // 试卷编号
	Answers       string `gorm:"column:answers" json:"answers"`                                   // 学生答案
	DetailedScore string `gorm:"detailed_score" json:"detailed_score"`                            //试卷题目详细分数
	TotalScore    string `gorm:"total_score" json:"total_score"`                                  // 总分数
	Status        int    `gorm:"status" json:"status"`                                            //状态
	ProblemStatus string `gorm:"problem_status" json:"problem_status"`                            // 题目状态
	Comment       string `gorm:"column:comment" form:"comment" binding:"required" json:"comment"` // 备注
	Times         int    `gorm:"times" json:"times"`                                              // 做的次数
}

type ReviewClass struct {
	ClassID string `gorm:"class_id" json:"class_id"`
	BaseClassInfo
	Count int `gorm:"count" json:"count"`
}

type ReviewStudent struct {
	StudentID  string    `gorm:"student_id" json:"student_id"`
	Name       string    `gorm:"name" json:"name"`
	Score      string    `gorm:"score" json:"score"`
	Status     int       `gorm:"status" json:"status"`
	UpdateTime time.Time `gorm:"update_time" json:"update_time"`
}

type ReviewUpdate struct {
	ReviewUpdateReq
	UpdateTime string `gorm:"column:update_time" json:"update_time"` // 更新时间
}

type ReviewUpdateReq struct {
	ExamID        string `gorm:"column:exam_id" json:"exam_id"`        // 试卷编号
	StudentID     string `gorm:"column:student_id" json:"student_id"`  //学生编号
	Answers       string `gorm:"column:answers" json:"answers"`        // 学生答案
	DetailedScore string `gorm:"detailed_score" json:"detailed_score"` //试卷题目详细分数
	TotalScore    string `gorm:"total_score" json:"total_score"`       // 总分数
	ProblemStatus string `gorm:"problem_status" json:"problem_status"` // 题目状态
	Status        int    `gorm:"status" json:"status"`                 //状态
	Times         int    `gorm:"times" json:"times"`                   // 做的次数
}

type ReviewStudentListReq struct {
	ClassID string `json:"class_id" form:"class_id"`
	ExamID  string `json:"exam_id" form:"exam_id"`
}

type ReviewStudentReq struct {
	ExamID    string `json:"exam_id" form:"exam_id"`
	StudentID string `json:"student_id" form:"student_id"`
}
