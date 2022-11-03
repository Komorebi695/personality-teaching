package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
	"time"
)

type TQuestion struct {
	ID         int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	QuestionId string    `json:"question_id" gorm:"column:question_id" description:"题目Id"`
	Name       string    `json:"name" gorm:"column:name" description:"题目名称"`
	Level      int       `json:"level" gorm:"column:level" description:"困难程度，1：容易，2：中等，3：困难"`
	Type       int       `json:"type" gorm:"column:type" description:"题目类型，1：选择题，2：填空题，3：问答题"`
	Context    string    `json:"context" gorm:"column:context" description:"题目内容"`
	Answer     string    `json:"answer" gorm:"column:answer" description:"题目答案"`
	CreateUser string    `json:"create_user" gorm:"column:create_user" description:"录入者"`
	UpdatedAt  time.Time `json:"update_time" gorm:"column:update_time;default:null" description:"修改时间"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_time;default:null" description:"创建时间"`
	IsDelete   int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (t *TQuestion) TableName() string {
	return "t_question"
}

func (t *TQuestion) Find(c *gin.Context, tx *gorm.DB) (*TQuestion, error) {
	out := &TQuestion{}
	err := tx.WithContext(c).Where("question_id = ?", t.QuestionId).Or("name = ?", t.Name).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *TQuestion) FindById(c *gin.Context, tx *gorm.DB) (*TQuestion, error) {
	out := &TQuestion{}
	err := tx.WithContext(c).Where("question_id = ?", t.QuestionId).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *TQuestion) PageList(c *gin.Context, tx *gorm.DB, param *model.QuestionListInput) ([]TQuestion, int64, error) {
	total := int64(0)
	var list []TQuestion
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.WithContext(c).Table(t.TableName()).Where("is_delete = 0")
	if param.Info != "" {
		query = query.Where("(name like ?)", "%"+param.Info+"%")
	}

	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (t *TQuestion) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}
