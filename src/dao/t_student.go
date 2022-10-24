package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type TStudent struct {
	ID          int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	StudentId   string `json:"student_id" gorm:"column:student_id" description:"学生编号,学号"`
	Name        string `json:"name" gorm:"column:name" description:"学生姓名"`
	Password    string `json:"password" gorm:"column:password" description:"密码, 加盐哈希"`
	College     string `json:"college" gorm:"column:college" description:"学院名称"`
	Major       string `json:"major" gorm:"column:major" description:"专业"`
	ClassId     string `json:"class_id" gorm:"column:class_id" description:"班级编号,UUID"`
	PhoneNumber string `json:"phone_number" gorm:"column:phone_number" description:"电话号码"`
	IsDelete    int8   `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (t *TStudent) TableName() string {
	return "t_student"
}

func (t *TStudent) Find(c *gin.Context, tx *gorm.DB, search *TStudent) (*TStudent, error) {
	out := &TStudent{}
	err := tx.WithContext(c).Where("student_id = ?", search.StudentId).Or("name = ?", search.Name).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *TStudent) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}

func (t *TStudent) PageList(c *gin.Context, tx *gorm.DB, param *model.StudentListInput) ([]TStudent, int64, error) {
	total := int64(0)
	var list []TStudent
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.WithContext(c).Table(t.TableName()).Where("is_delete = 0")
	if param.Info != "" {
		query = query.Where("(name like ? or student_id like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}

	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
