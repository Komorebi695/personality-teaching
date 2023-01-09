package model

// KnowledgeConnection 知识点联系表
type KnowledgeConnection struct {
	Id     int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	KnpId  string `json:"knp_id" gorm:"column:knp_id" description:"知识点编号"`
	PKnpId string `json:"p_knp_id" gorm:"column:p_knp_id" description:"前驱知识点编号"`
}
