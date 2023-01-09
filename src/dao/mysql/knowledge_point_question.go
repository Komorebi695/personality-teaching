package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type knowledgePointQuestionFunc interface {
	Find(c *gin.Context, tx *gorm.DB, questionId string) ([]*model.KnowledgePointQuestion, error)
	Save(c *gin.Context, tx *gorm.DB, param *model.KnowledgePointQuestion) error
	Delete(c *gin.Context, tx *gorm.DB, id int64) error
	DeleteAllById(c *gin.Context, tx *gorm.DB, knpId string) error
}

type KnowledgePointQuestionMySQL struct{}

var _ knowledgePointQuestionFunc = &KnowledgePointQuestionMySQL{}

func NewKnowledgePointQuestionMySQL() *KnowledgePointQuestionMySQL {
	return &KnowledgePointQuestionMySQL{}
}

func (t *KnowledgePointQuestionMySQL) TableName() string {
	return "t_knowledge_point_question"
}

// Find 根据questionId查找对应knpId列表
func (t *KnowledgePointQuestionMySQL) Find(c *gin.Context, tx *gorm.DB, questionId string) ([]*model.KnowledgePointQuestion, error) {
	var list []*model.KnowledgePointQuestion
	err := tx.WithContext(c).Table(t.TableName()).Where("question_id = ?", questionId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Save 保存问题知识点联系
func (t *KnowledgePointQuestionMySQL) Save(c *gin.Context, tx *gorm.DB, param *model.KnowledgePointQuestion) error {
	if err := tx.WithContext(c).Table(t.TableName()).Save(param).Error; err != nil {
		return err
	}
	return nil
}

// Delete 根据主键id删除问题知识点联系
func (t *KnowledgePointQuestionMySQL) Delete(c *gin.Context, tx *gorm.DB, id int64) error {
	if err := tx.WithContext(c).Table(t.TableName()).Delete(&model.KnowledgePointQuestion{}, id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAllById 批量刪除一个question下的所有知识点联系数据
func (t *KnowledgePointQuestionMySQL) DeleteAllById(c *gin.Context, tx *gorm.DB, knpId string) error {
	if err := tx.WithContext(c).Table(t.TableName()).Where("question_id = ?", knpId).Delete(t).Error; err != nil {
		return err
	}
	return nil
}
