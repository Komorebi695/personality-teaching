package model

type BaseStudentInfo struct {
	Name        string ` json:"name" form:"name" binding:"required" gorm:"column:name"`
	StudentNo   string ` json:"student_no" form:"student_no" binding:"required" gorm:"column:student_no"`
	College     string `json:"college" form:"college" binding:"required" gorm:"column:college"`
	Major       string `json:"major" form:"major" binding:"required" gorm:"column:major"`
	PhoneNumber string `json:"phone_number" form:"phone_number" gorm:"column:phone_number"`
}

type Student struct {
	StudentID string `gorm:"column:student_id"`
	Password  string `gorm:"column:password"`
	ClassID   string `gorm:"column:class_id"`
	BaseStudentInfo
}

type StudentQuestion struct {
	KnpID    string  `gorm:"column:knp_id"`
	AllScore float32 `gorm:"column:allscore"`
	Score    float32 `gorm:"column:score"`
	Answer   string  `gorm:"column:answer"`
}

type SearchStudentIDResp struct {
	StudentID string `json:"student_id" form:"student_id" binding:"required" gorm:"column:student_id"`
}

func (s Student) TableName() string {
	return "t_student"
}

type CreateStudentReq struct {
	BaseStudentInfo
}

type CreateStudentResp struct {
	StudentID string `json:"student_id" form:"student_id" binding:"required" gorm:"column:student_id"`
	BaseStudentInfo
}

type AddStudentToClassReq struct {
	StudentID string `form:"student_id" binding:"required"`
	ClassID   string `form:"class_id" binding:"required"`
}

type AddStudentClassResp struct {
	StudentID string `json:"student_id"`
	ClassID   string `json:"class_id"`
	BaseStudentInfo
}

type ClassStudentListReq struct {
	ClassID  string `form:"class_id" binding:"required"`
	PageNum  int    `form:"page_num" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
}

type ClassStudentListResp struct {
	StudentID string `json:"student_id" gorm:"column:student_id"`
	BaseStudentInfo
}

type EmptyClassStudentReq struct {
	PageNum  int    `form:"page_num" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Content  string `form:"content"`
}

type DeleteClassStudentReq struct {
	ClassID   string `form:"class_id" binding:"required"`
	StudentID string `form:"student_id" binding:"required"`
}

type SearchStudentReq struct {
	SearchText string `json:"search_text" form:"search_text" binding:"required"`
}

type DeleteStudentReq struct {
	StudentID string `form:"student_id" binding:"required"`
}
