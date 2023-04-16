package model

type Exam struct {
	ExamID          string `gorm:"column:exam_id" json:"exam_id"`                // 试卷编号
	CreateTeacherID string `gorm:"teacher_id;not null" json:"create_teacher_id"` // 创建老师编号
	BaseExamInfo
	UpdateTime string `gorm:"column:update_time" json:"update_time"` // 更新时间
	CreateTime string `gorm:"column:create_time" json:"create_time"` // 创建时间
}

func (e Exam) TableName() string {
	return "t_exam"
}

// BaseExamInfo ,试卷基本信息
type BaseExamInfo struct {
	ExamName  string `gorm:"column:exam_name" form:"exam_name" binding:"required" json:"exam_name"` // 试卷名称
	Questions string `gorm:"column:questions" form:"questions" binding:"required" json:"questions"` // 试题
	Comment   string `gorm:"column:comment" form:"comment" binding:"required" json:"comment"`       // 备注
}

// ExamAddReq ,新增试卷请求参数结构
type ExamAddReq struct {
	BaseExamInfo // 试卷基本信息
}

// ExamDeleteReq ,删除试卷请求参数结构
type ExamDeleteReq struct {
	ExamID string `gorm:"column:exam_id" form:"exam_id" binding:"required" json:"exam_id"` // 试卷编号
}

// ExamUpdateReq ,更新试卷请求参数结构
type ExamUpdateReq struct {
	ExamID string `gorm:"column:exam_id" form:"exam_id" binding:"required" json:"exam_id"` // 试卷编号
	BaseExamInfo
}

// ExamIDReq ,试卷详情
type ExamIDReq struct {
	ExamID string `gorm:"column:exam_id" form:"exam_id" binding:"required" json:"exam_id"`
}

// PagingReq ,分页参数
type PagingReq struct {
	Page     int `json:"page" form:"page"`           // 页数（第几页）
	PageSize int `json:"page_size" form:"page_size"` // 页面大小
}

type ExamListResp struct {
	Total    int        `gorm:"total" json:"total"`
	ExamList []ExamResp `gorm:"exam_list" json:"exam_list"`
}

// ExamResp ,试卷列表返回结构
type ExamResp struct {
	ExamID     string `gorm:"column:exam_id" form:"exam_id" binding:"required" json:"exam_id"`       // 试卷编号
	ExamName   string `gorm:"column:exam_name" form:"exam_name" binding:"required" json:"exam_name"` // 试卷名称
	Comment    string `gorm:"column:comment" form:"comment" binding:"required" json:"comment"`       // 备注
	UpdateTime string `gorm:"column:update_time" json:"update_time"`                                 // 更新时间
}

// ExamDetailResp ,试卷详情返回结构
type ExamDetailResp struct {
	ExamID string `gorm:"column:exam_id" form:"exam_id" binding:"required" json:"exam_id"` // 试卷编号
	BaseExamInfo
	UpdateTime string `gorm:"column:update_time" json:"update_time"` // 更新时间
}

type SearchReq struct {
	Text string `json:"text"`
}

type ReleaseStudentResp struct {
	StudentID string `json:"student_id" gorm:"column:student_id"`
	Name      string `json:"name"  gorm:"column:name"`
	College   string `json:"college"  gorm:"column:college"`
	Major     string `json:"major" gorm:"column:major"`
	Status    string `json:"status" gorm:"column:status"`
}

type ReleaseExamReq struct {
	ClassID
	ExamIDReq
}

type GetTeacherExamListReq struct {
	StudentID
}

type PostStudentExamAnswerReq struct {
	StudentID
	ExamIDReq
	Answer string `json:"answers" gorm:"column:answers"`
}

type StudentAllExams struct {
	ID int
	ExamDetailResp
}

type StudentReviewExams struct {
	StudentID
	ExamIDReq
	Status string `json:"status" gorm:"column:status" form:"status"`
}