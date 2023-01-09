package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type knowledgeConnectionFunc interface {
	Find(c *gin.Context, tx *gorm.DB, knpId string) ([]*model.KnowledgeConnection, error)
	Save(c *gin.Context, tx *gorm.DB, param *model.KnowledgeConnection) error
	Delete(c *gin.Context, tx *gorm.DB, id int64) error
	DeleteById(c *gin.Context, tx *gorm.DB, knpId string) error
	QueryNameById(c *gin.Context, tx *gorm.DB, knpId string) ([]*model.KnowledgeConnectionItem, error)
}

type KnowledgeConnectionMySQL struct{}

var _ knowledgeConnectionFunc = &KnowledgeConnectionMySQL{}

func NewKnowledgeConnectionMySQL() *KnowledgeConnectionMySQL {
	return &KnowledgeConnectionMySQL{}
}

func (t *KnowledgeConnectionMySQL) TableName() string {
	return "t_knowledge_connection"
}

// Find 根据knpId查找知识点连接列表
func (t *KnowledgeConnectionMySQL) Find(c *gin.Context, tx *gorm.DB, knpId string) ([]*model.KnowledgeConnection, error) {
	var list []*model.KnowledgeConnection
	err := tx.WithContext(c).Table(t.TableName()).Select("id, knp_id, p_knp_id").Where("knp_id = ?", knpId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Save 保存知识点连接
func (t *KnowledgeConnectionMySQL) Save(c *gin.Context, tx *gorm.DB, param *model.KnowledgeConnection) error {
	if err := tx.WithContext(c).Table(t.TableName()).Save(param).Error; err != nil {
		return err
	}
	return nil
}

// Delete 根据主键id删除知识点连接
func (t *KnowledgeConnectionMySQL) Delete(c *gin.Context, tx *gorm.DB, id int64) error {
	if err := tx.WithContext(c).Table(t.TableName()).Delete(&model.KnowledgeConnection{}, id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteById 批量刪除一个knowledge下的所有连接数据
func (t *KnowledgeConnectionMySQL) DeleteById(c *gin.Context, tx *gorm.DB, knpId string) error {
	if err := tx.WithContext(c).Table(t.TableName()).Where("knp_id = ?", knpId).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// QueryNameById 返回知识点连接的前驱知识点列表
func (t *KnowledgeConnectionMySQL) QueryNameById(c *gin.Context, tx *gorm.DB, knpId string) ([]*model.KnowledgeConnectionItem, error) {
	var list []*model.KnowledgeConnectionItem
	if err := tx.WithContext(c).Table(t.TableName()).Select("id, knp_id, p_knp_id").Table("t_knowledge_point A ,t_knowledge_connection B").Select("A.name p_name", "B.p_knp_id").Where("A.knp_id = B.p_knp_id and B.knp_id = ?", knpId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
