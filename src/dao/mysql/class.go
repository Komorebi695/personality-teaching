package mysql

import (
	"personality-teaching/src/model"

	"gorm.io/gorm"
)

const (
	classValid   int8 = 1 //合法记录
	classInvalid int8 = 0 //不合法记录
	classRows    int  = 10
)

type classFunc interface {
	InsertClass(teacherID string, c model.Class) error
	UpdateClass(m model.Class) error
	DeleteClass(teacherID string, classID string) error
	QueryClass(classID string) (model.Class, error)
	QueryClassList(teacherID string, pn int) ([]model.Class, error)
	CheckTeacherClass(teacherID string, classID string) (bool, error)
}

type ClassMySQL struct{}

var _ classFunc = &ClassMySQL{}

func NewClassMysql() *ClassMySQL {
	return &ClassMySQL{}
}

func (c *ClassMySQL) InsertClass(teacherID string, m model.Class) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("insert into `t_class`(`class_id`,`name`,`college`,`major`) values (?,?,?,?)",
			m.ClassID, m.Name, m.College, m.Major).Error; err != nil {
			return err
		}
		if err := tx.Exec("insert into `t_teacher_class`(`class_id`,`teacher_id`,`is_valid`) values (?,?,?)",
			m.ClassID, teacherID, classValid).Error; err != nil {
			return err
		}
		return nil
	})

}

func (c *ClassMySQL) UpdateClass(m model.Class) error {
	return db.Exec("update `t_class` set `name` = ?,`college` = ?,`major` = ? where class_id = ?",
		m.Name, m.College, m.Major, m.ClassID).Error
}

func (c *ClassMySQL) DeleteClass(teacherID string, classID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("delete from `t_class` where `class_id` = ?", classID).Error; err != nil {
			return err
		}
		if err := tx.Exec("update `t_teacher_class` set `is_valid` = ? where `teacher_id` = ? and `class_id` = ?",
			classInvalid, teacherID, classID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (c *ClassMySQL) QueryClass(classID string) (model.Class, error) {
	var m model.Class
	if err := db.Raw("select `class_id`,`name`,`college`,`major` from `t_class` where class_id = ?", classID).Scan(&m).Error; err != nil {
		return model.Class{}, err
	}
	return m, nil
}

func (c *ClassMySQL) QueryClassList(teacherID string, pn int) ([]model.Class, error) {
	var classes []model.Class
	offset := (pn - 1) * classRows
	count := classRows
	err := db.Raw("select `t_class`.`class_id`,`name`,`college`,`major` from `t_class` inner join `t_teacher_class` "+
		"on `t_class`.class_id = `t_teacher_class`.class_id "+
		"where teacher_id = ? and `is_valid` = ? limit ?,?", teacherID, classValid, offset, count).Scan(&classes).Error
	if err != nil {
		return []model.Class{}, err
	}
	return classes, nil
}

// CheckTeacherClass 有此数据返回true
func (c *ClassMySQL) CheckTeacherClass(teacherID string, classID string) (bool, error) {
	id := ""
	err := db.Raw("select `id` from `t_teacher_class` where `teacher_id` = ? and `class_id` = ? and `is_valid` = ?", teacherID, classID, classValid).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return !(id == ""), nil
}
