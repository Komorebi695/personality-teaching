package model

import "strings"

// QuestionListInput 题目列表输入
type QuestionListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" binding:""`
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
	QuestionId   string `json:"question_id" form:"question_id" comment:"question_id"`
	Type         int    `json:"type" form:"type" comment:"题目类型，1：选择题，2：填空题，3：问答题"`
	Level        int    `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	QuestionName string `json:"question_name" form:"question_name" comment:"题目名称,即题目主体"`
	Context      string `json:"context" form:"context" comment:"题目内容，存放题目中除了题干文字之外的内容，如图片地址"`
	Answer       string `json:"answer" form:"answer" comment:"题目答案"`
	CreateUser   string `json:"create_user" form:"create_user" comment:"录入者"`
}

type QuestionDeleteInput struct {
	QuestionId string `json:"question_id" form:"question_id" comment:"问题ID" binding:"required"`
}

// QuestionAddInput 问题添加输入
type QuestionAddInput struct {
	QuestionName       string                   `json:"question_name" form:"question_name" comment:"题目名称" binding:"required"`
	Context            string                   `json:"context" form:"context" comment:"题目内容"`
	Answer             string                   `json:"answer" form:"answer" comment:"题目答案"`
	Type               int                      `json:"type" form:"type" comment:"题目类型，1：选择题，2：填空题，3：问答题" binding:"required"`
	Level              int                      `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难" binding:"required"`
	CreateUser         string                   `json:"create_user" form:"create_user" comment:"录入者" binding:"required"`
	QuestionOptionList []QuestionOptionAddInput `json:"question_option_list" form:"question_option_list" comment:"选项信息，非选择题则为空"`
	KnpId              string                   `json:"knp_id" form:"column:knp_id" comment:"上级知识点编号，以逗号分隔"`
}

func (t *QuestionAddInput) GetKnpIdByModel() []string {
	return strings.Split(t.KnpId, ",")
}

type QuestionDetailInput struct {
	QuestionId string `json:"question_id" form:"question_id" comment:"题目ID" binding:"required"`
}

// QuestionUpdateInput 题目修改输入
type QuestionUpdateInput struct {
	QuestionId   string                   `json:"question_id" form:"question_id" comment:"问题id" binding:"required"`
	QuestionName string                   `json:"question_name" form:"question_name" comment:"题目名称" binding:"required"`
	Context      string                   `json:"context" form:"context" comment:"题目内容"`
	Option       []QuestionOptionAddInput `json:"question_option_list" form:"question_option_list" comment:"选项信息，非选择题则为空"`
	Answer       string                   `json:"answer" form:"answer" comment:"题目答案"`
	Type         int                      `json:"type" form:"type" comment:"题目类型，1：选择题，2：填空题，3：问答题" binding:"required"`
	Level        int                      `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难" binding:"required"`
	CreateUser   string                   `json:"create_user" form:"create_user" comment:"录入者" binding:"required"`
	KnpId        string                   `json:"knp_id" form:"column:knp_id" comment:"上级知识点编号,以逗号分隔"`
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

// QuestionOptionAddInput 选项添加输入
type QuestionOptionAddInput struct {
	Context  string `json:"context" form:"context" comment:"选项内容" binding:"required"`
	IsAnswer int8   `json:"is_answer" form:"is_answer" comment:"是否为正确答案，0：否，1：是" binding:"required"`
}
