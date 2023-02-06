package model

import "time"

// KnowledgePointQuestion 知识点题目联系表
type KnowledgePointQuestion struct {
	Id         int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	KnpId      string    `json:"knp_id" gorm:"column:knp_id" description:"知识点编号编号"`
	QuestionId string    `json:"question_id" gorm:"column:question_id" description:"题目编号"`
	UpdatedAt  time.Time `json:"update_time" gorm:"column:update_time;default:null" description:"修改时间"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_time;default:null" description:"创建时间"`
}
