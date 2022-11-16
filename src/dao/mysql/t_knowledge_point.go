package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
	"time"
)

type TKnowledgePoint struct {
	Id          int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	KnpId       string    `json:"knp_id" gorm:"column:knp_id" description:"知识点编号"`
	Name        string    `json:"name" gorm:"column:name" description:"知识点名称"`
	ParentKnpId string    `json:"parent_knp_id" gorm:"column:parent_knp_id" description:"父知识点编号，若没有，则是它自身"`
	Level       int       `json:"level" gorm:"column:level" description:"困难程度，1：容易，2：中等，3：困难"`
	Context     string    `json:"context" gorm:"column:context" description:"知识点内容"`
	CreateUser  string    `json:"create_user" gorm:"column:create_user" description:"录入者"`
	UpdatedAt   time.Time `json:"update_time" gorm:"column:update_time" description:"修改时间"`
	CreatedAt   time.Time `json:"create_time" gorm:"column:create_time" description:"创建时间"`
}

func (t *TKnowledgePoint) TableName() string {
	return "t_knowledge_point"
}

func (t *TKnowledgePoint) FindOneById(c *gin.Context, tx *gorm.DB) (*TKnowledgePoint, error) {
	point := &TKnowledgePoint{}
	err := tx.WithContext(c).Where("knp_id = ?", t.KnpId).First(point).Error
	return point, err
}

func (t *TKnowledgePoint) FindByName(c *gin.Context, tx *gorm.DB) (*TKnowledgePoint, error) {
	point := &TKnowledgePoint{}
	err := tx.WithContext(c).Where("name = ?", t.Name).First(point).Error
	return point, err
}

func (t *TKnowledgePoint) PageList(c *gin.Context, tx *gorm.DB, param *model.KnowledgePointListInput) ([]TKnowledgePoint, int64, error) {
	total := int64(0)
	var list []TKnowledgePoint
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.WithContext(c).Table(t.TableName())
	if param.Info != "" {
		query = query.Where("(name like ?)", "%"+param.Info+"%")
	}

	query.Limit(param.PageSize).Offset(-1).Count(&total)
	if err := query.Limit(param.PageSize).Offset(offset).Order("id asc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return list, total, nil
}
func (t *TKnowledgePoint) PageListOneStage(c *gin.Context, tx *gorm.DB) ([]model.KnpOneStageListItemOutput, error) {
	var list []model.KnpOneStageListItemOutput

	if err := tx.WithContext(c).Table("t_knowledge_point A").Select("A.knp_id,A.name,A.context,A.parent_knp_id,A.level,B.knp_id,B.level,B.name,B.context,B.parent_knp_id").Joins("inner join t_knowledge_point B on A.knp_id = B.parent_knp_id and B.knp_id = A.parent_knp_id").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}
func (t *TKnowledgePoint) Delete(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Delete(t).Error; err != nil {
		return err
	}
	return nil
}
func (t *TKnowledgePoint) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}

//FindKnowledgeChildren 查找子节点
func (t *TKnowledgePoint) FindKnowledgeChildren(c *gin.Context, tx *gorm.DB) ([]*TKnowledgePoint, error) {
	var list []*TKnowledgePoint
	err := tx.WithContext(c).Table(t.TableName()).Where("parent_knp_id = ?", t.KnpId).Not("knp_id = ?", t.KnpId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
