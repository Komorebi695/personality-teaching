package model

import (
	"strings"
	"time"
)

// Question 题目表
type Question struct {
	ID         int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	QuestionId string `json:"question_id" gorm:"column:question_id" description:"题目Id"`
	QuestionBase
	CreateUser string    `json:"create_user" gorm:"column:create_user" description:"录入者"`
	UpdatedAt  time.Time `json:"update_time" gorm:"column:update_time;default:null" description:"修改时间"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_time;default:null" description:"创建时间"`
}

// QuestionBase 题目信息表
type QuestionBase struct {
	Name    string `json:"name" gorm:"column:name" description:"题目名称"`
	Level   int    `json:"level" gorm:"column:level" description:"困难程度，1：容易，2：中等，3：困难"`
	Type    int    `json:"type" gorm:"column:type" description:"题目类型，1：单选题，2：多选题，3：填空题，4：问答题"`
	Context string `json:"context" gorm:"column:context" description:"题目内容"`
	Answer  string `json:"answer" gorm:"column:answer" description:"题目答案"`
}

// QuestionInfo 题目信息表
type QuestionInfo struct {
	Question
	AnswerContext  string            `json:"answer_context" form:"answer_context" comment:"答案解析"`
	QuestionOption []*QuestionOption `json:"question_option" comment:"选项信息"`
}

// QuestionDetail 题目详细
type QuestionDetail struct {
	QuestionInfo
	KnowledgePointList []*KnowledgePoint `json:"knowledge_point_list" description:"知识点列表"`
}

// QuestionListInput 题目列表输入
type QuestionListInput struct {
	Context  string `json:"context" form:"context" comment:"关键词"`
	Type     int    `json:"type" form:"type" comment:"题目类型，1：单选题，2：多选题，3：填空题，4：问答题"`
	Level    int    `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	KnpId    string `json:"knp_id" form:"knp_id" comment:"知识点Id"`
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
	QuestionBase
	QuestionId    string            `json:"question_id" gorm:"column:question_id" description:"题目Id"`
	Option        []*QuestionOption `json:"question_option_list" form:"question_option_list" comment:"选项信息，非选择题则为空"`
	AnswerContext string            `json:"answer_context" form:"answer_context" comment:"答案解析"`
	CreateUser    string            `json:"create_user" form:"create_user" comment:"录入者"`
}

type QuestionDeleteInput struct {
	QuestionId string `json:"question_id" form:"question_id" comment:"问题ID" binding:"required"`
}

// QuestionAddInput 问题添加输入
type QuestionAddInput struct {
	QuestionBase
	AnswerContext      string            `json:"answer_context" form:"answer_context" comment:"答案解析"`
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

// QuestionDetailInput 题目详情输入
type QuestionDetailInput struct {
	QuestionId string `json:"question_id" form:"question_id" comment:"题目ID" binding:"required"`
}

// QuestionUpdateInput 题目修改输入
type QuestionUpdateInput struct {
	QuestionBase
	QuestionId    string            `json:"question_id" form:"question_id" comment:"问题id" binding:"required"`
	Option        []*QuestionOption `json:"question_option_list" form:"question_option_list" comment:"选项信息，非选择题则为空"`
	AnswerContext string            `json:"answer_context" form:"answer_context" comment:"答案解析"`
	CreateUser    string            `json:"create_user" form:"create_user" comment:"录入者" binding:"required"`
	KnpId         string            `json:"knp_id" form:"column:knp_id" comment:"上级知识点编号,以逗号分隔"`
}

func (t *QuestionUpdateInput) GetKnpIdByModel() []string {
	return strings.Split(t.KnpId, ",")
}

// QuestionOption 选项添加输入
type QuestionOption struct {
	Context string `json:"Context" form:"Context" comment:"选项内容" binding:"required"`
}
