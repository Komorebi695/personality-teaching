package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type TKnowledgeConnection struct {
	Id     int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	KnpId  string `json:"knp_id" gorm:"column:knp_id" description:"知识点编号"`
	PKnpId string `json:"p_knp_id" gorm:"column:p_knp_id" description:"前驱知识点编号"`
}

func (t *TKnowledgeConnection) TableName() string {
	return "t_knowledge_connection"
}

func (t *TKnowledgeConnection) Find(c *gin.Context, tx *gorm.DB) ([]*TKnowledgeConnection, error) {
	var list []*TKnowledgeConnection
	err := tx.WithContext(c).Where("knp_id = ?", t.KnpId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (t *TKnowledgeConnection) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *TKnowledgeConnection) Delete(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// DeleteById 批量刪除一个knowledge下的所有关联数据
func (t *TKnowledgeConnection) DeleteById(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Where("knp_id = ?", t.KnpId).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// QueryNameById 返回知识点连接的前驱知识点列表
func (t *TKnowledgeConnection) QueryNameById(c *gin.Context, tx *gorm.DB) ([]*model.KnowledgeConnectionItem, error) {
	var list []*model.KnowledgeConnectionItem
	if err := tx.WithContext(c).Table("t_knowledge_point A ,t_knowledge_connection B").Select("A.name p_name", "B.p_knp_id").Where("A.knp_id = B.p_knp_id and B.knp_id = ?", t.KnpId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
