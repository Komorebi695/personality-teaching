package model

import (
	"strings"
)

// KnowledgePointListInput 知识点列表输入
type KnowledgePointListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" binding:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" binding:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" binding:"required"` //每页条数
}

// KnowledgePointListOutput 知识点列表输出
type KnowledgePointListOutput struct {
	Total int64                          `json:"total" form:"total" comment:"总数"`  //总数
	List  []KnowledgePointListItemOutput `json:"list" form:"list" comment:"知识点列表"` //列表
}

// KnowledgePointListItemOutput 知识点列表输出的主体
type KnowledgePointListItemOutput struct {
	KnpId       string `json:"knp_id" form:"knp_id" comment:"knp_id,知识点编号"`
	ParentKnpId string `json:"parent_knp_id" form:"parent_knp_id" comment:"父知识点编号，若没有，则是它自身"`
	Name        string `json:"name" form:"name" comment:"知识点名称"`
	Level       int    `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	Context     string `json:"context" form:"context" comment:"知识点内容"`
	CreateUser  string `json:"create_user" form:"create_user" comment:"录入者"`
}

// KnowledgePointOneStageListOutput 知识点一级列表输出
type KnowledgePointOneStageListOutput struct {
	List []KnpOneStageListItemOutput `json:"list" form:"list" comment:"知识点一级列表"` //列表
}

// KnpOneStageListItemOutput 知识点一级列表输出的主体
type KnpOneStageListItemOutput struct {
	KnpId       string `json:"knp_id" form:"knp_id" comment:"knp_id,知识点编号"`
	ParentKnpId string `json:"parent_knp_id" form:"parent_knp_id" comment:"父知识点编号，若没有，则是它自身"`
	Name        string `json:"name" form:"name" comment:"知识点名称"`
	Level       int    `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	Context     string `json:"context" form:"context" comment:"知识点内容"`
}
type KnowledgePointDeleteInput struct {
	KnpId string `json:"knp_id" form:"knp_id" comment:"知识点ID" binding:"required"` //知识点ID
}

// KnowledgePointAddInput 知识点添加输入
type KnowledgePointAddInput struct {
	Name        string `json:"name" form:"name" comment:"知识点名称"`
	ParentKnpId string `json:"parent_knp_id" form:"parent_knp_id" comment:"父知识点编号，若没有，则是它自身"`
	Level       int    `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	Context     string `json:"context" form:"context" comment:"知识点内容"`
	CreateUser  string `json:"create_user" form:"create_user" comment:"录入者"`
}

type KnowledgePointDetailInput struct {
	KnpId string `json:"knp_id" form:"knp_id" comment:"知识点编号" binding:"required"`
}

// KnowledgePointUpdateInput 知识点修改输入
type KnowledgePointUpdateInput struct {
	KnpId       string `json:"knp_id" form:"knp_id" comment:"knp_id,知识点编号"`
	Name        string `json:"name" form:"name" comment:"知识点名称"`
	ParentKnpId string `json:"parent_knp_id" form:"parent_knp_id" comment:"父知识点编号，若没有，则是它自身"`
	Level       int    `json:"level" form:"level" comment:"困难程度，1：容易，2：中等，3：困难"`
	Context     string `json:"context" form:"context" comment:"知识点内容"`
	CreateUser  string `json:"create_user" form:"create_user" comment:"录入者"`
}

// KnpConnectionUpdateInput 知识点联系修改输入
type KnpConnectionUpdateInput struct {
	KnpId  string `json:"knp_id" form:"knp_id" comment:"知识点编号"`
	PKnpId string `json:"p_knp_id" form:"p_knp_id" comment:"前驱知识点编号，多个编号间用,隔开"`
}

func (t *KnpConnectionUpdateInput) GetKnpIdByModel() []string {
	return strings.Split(t.PKnpId, ",")
}
