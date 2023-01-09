package model

import (
	"strings"
	"time"
)

// Question 题目表
type Question struct {
	ID         int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	QuestionId string    `json:"question_id" gorm:"column:question_id" description:"题目Id"`
	Name       string    `json:"name" gorm:"column:name" description:"题目名称"`
	Level      int       `json:"level" gorm:"column:level" description:"困难程度，1：容易，2：中等，3：困难"`
	Type       int       `json:"type" gorm:"column:type" description:"题目类型，1：选择题，2：填空题，3：问答题"`
	Context    string    `json:"context" gorm:"column:context" description:"题目内容"`
	Answer     string    `json:"answer" gorm:"column:answer" description:"题目答案"`
	CreateUser string    `json:"create_user" gorm:"column:create_user" description:"录入者"`
	UpdatedAt  time.Time `json:"update_time" gorm:"column:update_time;default:null" description:"修改时间"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_time;default:null" description:"创建时间"`
}

// QuestionDetail 题目详细
type QuestionDetail struct {
	QuestionInfo       *Question         `json:"problem_info" description:"题目信息"`
	QuestionOption     []*QuestionOption `json:"question_option" comment:"选项信息"`
	KnowledgePointList []*KnowledgePoint `json:"knowledge_point_list" description:"知识点列表"`
}

// QuestionListInput 题目列表输入
type QuestionListInput struct {
	Context  string `json:"context" form:"context" comment:"关键词"`
	Type     int    `json:"type" form:"type" comment:"题目类型，1：选择题，2：填空题，3：问答题"`
	Level    int    `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" binding:"required"`
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" binding:"required"`
}

// QuestionListOutput 题目列表输出
type QuestionListOutput struct {
	Total int64                    `json:"total" form:"total" comment:"总数"`
	List  []QuestionListItemOutput `json:"list" form:"list" comment:"题目列表"`
}

// QuestionListItemOutput 题目列表输出的主体
type QuestionListItemOutput struct {
	QuestionId   string            `json:"question_id" form:"question_id" comment:"question_id"`
	Type         int               `json:"type" form:"type" comment:"题目类型，1：选择题，2：填空题，3：问答题"`
	Level        int               `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	QuestionName string            `json:"question_name" form:"question_name" comment:"题目名称,即题目主体"`
	Context      string            `json:"context" form:"context" comment:"题目内容"`
	Option       []*QuestionOption `json:"question_option_list" form:"question_option_list" comment:"选项信息，非选择题则为空"`
	Answer       string            `json:"answer" form:"answer" comment:"题目答案"`
	CreateUser   string            `json:"create_user" form:"create_user" comment:"录入者"`
}

type QuestionDeleteInput struct {
	QuestionId string `json:"question_id" form:"question_id" comment:"问题ID" binding:"required"`
}

// QuestionAddInput 问题添加输入
type QuestionAddInput struct {
	QuestionName       string            `json:"question_name" form:"question_name" comment:"题目名称" binding:"required"`
	Context            string            `json:"context" form:"context" comment:"题目内容"`
	Answer             string            `json:"answer" form:"answer" comment:"题目答案"`
	Type               int               `json:"type" form:"type" comment:"题目类型，1：选择题，2：填空题，3：问答题" binding:"required"`
	Level              int               `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难" binding:"required"`
	CreateUser         string            `json:"create_user" form:"create_user" comment:"录入者" binding:"required"`
	QuestionOptionList []*QuestionOption `json:"question_option_list" form:"question_option_list" comment:"选项信息，非选择题则为空"`
	KnpId              string            `json:"knp_id" form:"column:knp_id" comment:"上级知识点编号，以逗号分隔"`
}

// GetKnpIdByModel 知识点拆分
func (t *QuestionAddInput) GetKnpIdByModel() []string {
	if t.KnpId == "" {
		return nil
	}
	return strings.Split(t.KnpId, ",")
}

type QuestionDetailInput struct {
	QuestionId string `json:"question_id" form:"question_id" comment:"题目ID" binding:"required"`
}

// QuestionUpdateInput 题目修改输入
type QuestionUpdateInput struct {
	QuestionId   string            `json:"question_id" form:"question_id" comment:"问题id" binding:"required"`
	QuestionName string            `json:"question_name" form:"question_name" comment:"题目名称" binding:"required"`
	Context      string            `json:"context" form:"context" comment:"题目内容"`
	Option       []*QuestionOption `json:"question_option_list" form:"question_option_list" comment:"选项信息，非选择题则为空"`
	Answer       string            `json:"answer" form:"answer" comment:"题目答案"`
	Type         int               `json:"type" form:"type" comment:"题目类型，1：选择题，2：填空题，3：问答题" binding:"required"`
	Level        int               `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难" binding:"required"`
	CreateUser   string            `json:"create_user" form:"create_user" comment:"录入者" binding:"required"`
	KnpId        string            `json:"knp_id" form:"column:knp_id" comment:"上级知识点编号,以逗号分隔"`
}

func (t *QuestionUpdateInput) GetKnpIdByModel() []string {
	return strings.Split(t.KnpId, ",")
}

// QuestionOptionOutput 选项信息输出
type QuestionOptionOutput struct {
	ID       int64  `json:"id" form:"id" comment:"id"`
	Context  string `json:"context" form:"context" comment:"选项内容"`
	IsAnswer int8   `json:"is_answer" form:"is_answer" comment:"是否为正确答案，0：否，1：是"`
}

// QuestionOption 选项添加输入
type QuestionOption struct {
	Context string `json:"Context" form:"Context" comment:"选项内容" binding:"required"`
}
