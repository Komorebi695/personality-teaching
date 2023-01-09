package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type knowledgePointFunc interface {
	FindOneById(c *gin.Context, tx *gorm.DB, knpId string) (*model.KnowledgePoint, error)
	FindByName(c *gin.Context, tx *gorm.DB, name string) (*model.KnowledgePoint, error)
	PageList(c *gin.Context, tx *gorm.DB, param *model.KnowledgePointListInput) ([]model.KnowledgePointListItemOutput, int64, error)
	PageListOneStage(c *gin.Context, tx *gorm.DB) ([]model.KnpOneStageListItemOutput, error)
	Delete(c *gin.Context, tx *gorm.DB, id int64) error
	Save(c *gin.Context, tx *gorm.DB, param *model.KnowledgePoint) error
	FindKnowledgeChildren(c *gin.Context, tx *gorm.DB, knpId string) ([]*model.KnowledgePointBase, error)
}

type KnowledgePointMySQL struct{}

var _ knowledgePointFunc = &KnowledgePointMySQL{}

func NewKnowledgePointMySQL() *KnowledgePointMySQL {
	return &KnowledgePointMySQL{}
}

// TableName 表名
func (t *KnowledgePointMySQL) TableName() string {
	return "t_knowledge_point"
}

// FindOneById 根据Id查找一条知识点
func (t *KnowledgePointMySQL) FindOneById(c *gin.Context, tx *gorm.DB, knpId string) (*model.KnowledgePoint, error) {
	point := &model.KnowledgePoint{}
	err := tx.WithContext(c).Table(t.TableName()).Where("knp_id = ?", knpId).First(point).Error
	return point, err
}

// FindByName 根据name查找一条知识点
func (t *KnowledgePointMySQL) FindByName(c *gin.Context, tx *gorm.DB, name string) (*model.KnowledgePoint, error) {
	point := &model.KnowledgePoint{}
	err := tx.WithContext(c).Table(t.TableName()).Where("name = ?", name).First(point).Error
	return point, err
}

// PageList 查询知识点基本信息列表
func (t *KnowledgePointMySQL) PageList(c *gin.Context, tx *gorm.DB, param *model.KnowledgePointListInput) ([]model.KnowledgePointListItemOutput, int64, error) {
	total := int64(0)
	var list []model.KnowledgePointListItemOutput
	offset := (param.PageNo - 1) * param.PageSize
	query := tx.WithContext(c).Table(t.TableName()).Select("knp_id, name, parent_knp_id, level, context,create_user")
	if param.Info != "" {
		query = query.Where("(name like ?)", "%"+param.Info+"%")
	}
	query.Limit(param.PageSize).Offset(-1).Count(&total)
	if err := query.Limit(param.PageSize).Offset(offset).Order("id asc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return list, total, nil
}

// PageListOneStage 查询第一级知识点列表
func (t *KnowledgePointMySQL) PageListOneStage(c *gin.Context, tx *gorm.DB) ([]model.KnpOneStageListItemOutput, error) {
	var list []model.KnpOneStageListItemOutput
	err := tx.WithContext(c).
		Table(t.TableName() + " A").
		Select("A.knp_id,A.name,A.context,A.parent_knp_id,A.level,B.knp_id,B.level,B.name,B.context,B.parent_knp_id").
		Joins("inner join " + t.TableName() + " B on A.knp_id = B.parent_knp_id and B.knp_id = A.parent_knp_id").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}

// Delete 根据Id删除知识点
func (t *KnowledgePointMySQL) Delete(c *gin.Context, tx *gorm.DB, id int64) error {
	if err := tx.WithContext(c).Table(t.TableName()).Delete(&model.KnowledgePoint{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Save 知识点保存
func (t *KnowledgePointMySQL) Save(c *gin.Context, tx *gorm.DB, param *model.KnowledgePoint) error {
	return tx.WithContext(c).Table(t.TableName()).Save(param).Error
}

//FindKnowledgeChildren 根据KnpId查找子节点基本信息列表
func (t *KnowledgePointMySQL) FindKnowledgeChildren(c *gin.Context, tx *gorm.DB, knpId string) ([]*model.KnowledgePointBase, error) {
	var list []*model.KnowledgePointBase
	err := tx.WithContext(c).Table(t.TableName()).
		Select("id, knp_id, parent_knp_id").
		Where("parent_knp_id = ?", knpId).Not("knp_id = ?", knpId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
