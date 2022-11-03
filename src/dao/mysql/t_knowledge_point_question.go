package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type TKnowledgePointQuestion struct {
	Id         int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	KnpId      string    `json:"knp_id" gorm:"column:knp_id" description:"知识点编号编号"`
	QuestionId string    `json:"question_id" gorm:"column:question_id" description:"题目编号"`
	UpdatedAt  time.Time `json:"update_time" gorm:"column:update_time;default:null" description:"修改时间"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_time;default:null" description:"创建时间"`
}

func (t *TKnowledgePointQuestion) TableName() string {
	return "t_knowledge_point_question"
}

func (t *TKnowledgePointQuestion) Find(c *gin.Context, tx *gorm.DB, search *TKnowledgePointQuestion) ([]*TKnowledgePointQuestion, error) {
	var list []*TKnowledgePointQuestion
	err := tx.WithContext(c).Where("question_id = ?", search.QuestionId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (t *TKnowledgePointQuestion) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *TKnowledgePointQuestion) Delete(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// DeleteById 批量刪除一个question下的所有数据
func (t *TKnowledgePointQuestion) DeleteById(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Where("question_id = ?", t.QuestionId).Delete(t).Error; err != nil {
		return err
	}
	return nil
}
