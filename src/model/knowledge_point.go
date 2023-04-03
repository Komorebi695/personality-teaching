package model

import (
	"strings"
	"time"
)

// KnowledgePointBase 知识点基本编号
type KnowledgePointBase struct {
	KnpId       string `json:"knp_id" gorm:"column:knp_id" description:"知识点编号"`
	ParentKnpId string `json:"parent_knp_id" gorm:"column:parent_knp_id" description:"父知识点编号，若没有，则是它自身"`
}

// KnowledgePointInfo 知识点基本信息
type KnowledgePointInfo struct {
	Name    string `json:"name" gorm:"column:name" description:"知识点名称"`
	Level   int    `json:"level" gorm:"column:level" description:"困难程度，1：容易，2：中等，3：困难"`
	Context string `json:"context" gorm:"column:context" description:"知识点内容"`
}

// KnowledgePointImage 知识点图片储存
type KnowledgePointFile struct {
	Id     int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT" description:"自增主键"`
	CosUrl string `json:"cos_url" gorm:"column:Cos_url" description:"图片返回地址"`
}

// KnowledgePoint 知识点表
type KnowledgePoint struct {
	Id int64 `json:"id" gorm:"primary_key" description:"自增主键"`
	KnowledgePointBase
	KnowledgePointInfo
	CreateUser string    `json:"create_user" gorm:"column:create_user" description:"录入者"`
	UpdatedAt  time.Time `json:"update_time" gorm:"column:update_time" description:"修改时间"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_time" description:"创建时间"`
}

// KnowledgePointDetail 知识点详情表
type KnowledgePointDetail struct {
	Info                    *KnowledgePoint            `json:"info" description:"知识点信息"`
	Children                []*KnowledgePointBase      `json:"children" description:"子知识点列表信息"`
	KnowledgeConnectionList []*KnowledgeConnectionItem `json:"knowledge_connection_list" description:"知识点联系列表"`
}

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
	KnowledgePointBase
	KnowledgePointInfo
	CreateUser string `json:"create_user" form:"create_user" comment:"录入者"`
}

// KnpOneStageListOutput 知识点一级列表输出
type KnpOneStageListOutput struct {
	List []KnpOneStageListItemOutput `json:"list" form:"list" comment:"知识点一级列表"` //列表
}

// KnpOneStageListItemOutput 知识点一级列表输出的主体
type KnpOneStageListItemOutput struct {
	KnowledgePointBase
	KnowledgePointInfo
}

// KnowledgePointDeleteInput 知识点删除输入
type KnowledgePointDeleteInput struct {
	KnpId string `json:"knp_id" form:"knp_id" comment:"知识点ID" binding:"required"` //知识点ID
}

// KnowledgePointAddInput 知识点添加输入
type KnowledgePointAddInput struct {
	KnowledgePointInfo
	ParentKnpId string `json:"parent_knp_id" form:"parent_knp_id" comment:"父知识点编号，若没有，则是它自身"`
	CreateUser  string `json:"create_user" form:"create_user" comment:"录入者"`
}

// KnowledgePointDetailInput 知识点详情输入
type KnowledgePointDetailInput struct {
	KnpId string `json:"knp_id" form:"knp_id" comment:"知识点编号" binding:"required"`
}

// KnowledgeConnectionItem 知识点连接主体
type KnowledgeConnectionItem struct {
	PName  string `json:"p_name" form:"p_name" comment:"前驱知识点名字" `
	PKnpId string `json:"p_knp_id" form:"p_knp_id" comment:"前驱知识点编号"`
}

// KnowledgePointUpdateInput 知识点修改输入
type KnowledgePointUpdateInput struct {
	KnowledgePointBase
	KnowledgePointInfo
	CreateUser string `json:"create_user" form:"create_user" comment:"录入者"`
}

// KnpConnectionUpdateInput 知识点联系修改输入
type KnpConnectionUpdateInput struct {
	KnpId  string `json:"knp_id" form:"knp_id" comment:"知识点编号"`
	PKnpId string `json:"p_knp_id" form:"p_knp_id" comment:"前驱知识点编号，多个编号间用,隔开"`
}

// GetKnpIdByModel 知识点联系拆分
func (t *KnpConnectionUpdateInput) GetKnpIdByModel() []string {
	return strings.Split(t.PKnpId, ",")
}
