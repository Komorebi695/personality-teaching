package mysql

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type questionFunc interface {
	FindOnce(c *gin.Context, tx *gorm.DB, keyWord string) (*model.Question, error)
	PageList(c *gin.Context, tx *gorm.DB, param *model.QuestionListInput) ([]model.Question, int64, error)
	Save(c *gin.Context, tx *gorm.DB, param *model.Question) error
	Delete(c *gin.Context, tx *gorm.DB, id int64) error
}

type QuestionMySQL struct{}

var _ questionFunc = &QuestionMySQL{}

func NewQuestionMySQL() *QuestionMySQL {
	return &QuestionMySQL{}
}

func (q *QuestionMySQL) TableName() string {
	return "t_question"
}

// FindOnce 根据题目Id或题目名字查询
func (q *QuestionMySQL) FindOnce(c *gin.Context, tx *gorm.DB, keyWord string) (*model.Question, error) {
	out := &model.Question{}
	err := tx.WithContext(c).Table(q.TableName()).Where("question_id = ?", keyWord).Or("name = ?", keyWord).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (q *QuestionMySQL) PageList(c *gin.Context, tx *gorm.DB, param *model.QuestionListInput) ([]model.Question, int64, error) {
	total := int64(0)
	var list []model.Question
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.WithContext(c).Table(q.TableName())
	if param.Context != "" {
		query = query.Where("(context like ?)", "%"+param.Context+"%")
	}
	if param.Type != 0 {
		query = query.Where("(type = ?)", param.Type)
	}
	if param.Level != 0 {
		query = query.Where("(level = ?)", param.Level)
	}

	query.Limit(param.PageSize).Offset(-1).Count(&total)
	if err := query.Limit(param.PageSize).Offset(offset).Order("id asc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	return list, total, nil
}

func (q *QuestionMySQL) Save(c *gin.Context, tx *gorm.DB, param *model.Question) error {
	return tx.WithContext(c).Table(q.TableName()).Save(param).Error
}

// Delete 根据Id删除问题
func (q *QuestionMySQL) Delete(c *gin.Context, tx *gorm.DB, id int64) error {
	if err := tx.WithContext(c).Table(q.TableName()).Delete(&model.Question{}, id).Error; err != nil {
		return err
	}
	return nil
}
