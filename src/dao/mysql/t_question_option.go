package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type TQuestionOption struct {
	ID         int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	QuestionId string    `json:"question_id" gorm:"column:question_id" description:"题目id"`
	Context    string    `json:"context" gorm:"column:context" description:"选项内容"`
	IsAnswer   int8      `json:"is_answer" gorm:"column:is_answer" description:"是否为正确答案，0：否，1：是"`
	UpdatedAt  time.Time `json:"update_time" gorm:"column:update_time;default:null" description:"修改时间"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_time;default:null" description:"创建时间"`
}

func (t *TQuestionOption) TableName() string {
	return "t_question_option"
}

func (t *TQuestionOption) Find(c *gin.Context, tx *gorm.DB) (*TQuestionOption, error) {
	out := &TQuestionOption{}
	err := tx.WithContext(c).Where("id = ?", t.ID).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *TQuestionOption) FindByQuestionId(c *gin.Context, tx *gorm.DB) ([]*TQuestionOption, error) {
	var list []*TQuestionOption
	err := tx.WithContext(c).Where("question_id = ?", t.QuestionId).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (t *TQuestionOption) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}

func (t TQuestionOption) Delete(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Delete(t).Error; err != nil {
		return err
	}
	return nil
}

// DeleteById 批量刪除一个question下的所有数据
func (t *TQuestionOption) DeleteById(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Where("question_id = ?", t.QuestionId).Delete(t).Error; err != nil {
		return err
	}
	return nil
}
