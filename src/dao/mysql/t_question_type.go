package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TQuestionType struct {
	Id   int8   `json:"id" gorm:"primary_key" description:"自增主键"`
	Name string `json:"type_name" gorm:"column:type_name" description:"题目类型"`
}

func (t *TQuestionType) TableName() string {
	return "t_question_type"
}

func (t *TQuestionType) Find(c *gin.Context, tx *gorm.DB) (*TQuestionType, error) {
	point := &TQuestionType{}
	err := tx.WithContext(c).First(point).Error
	return point, err
}

func (t *TQuestionType) Delete(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Delete(t).Error; err != nil {
		return err
	}
	return nil
}
func (t *TQuestionType) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}
