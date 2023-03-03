package mysql

import (
	"personality-teaching/src/model"
)

type TeacherMySQL struct {
}

type teacherFunc interface {
	QueryAllByName(username string) (model.Teacher, error)

	QueryAllByID(teacherID string) (model.Teacher, error)
}

var _ teacherFunc = &TeacherMySQL{}

func NewTeacherMysql() *TeacherMySQL {
	return &TeacherMySQL{}
}

func (t *TeacherMySQL) QueryAllByName(username string) (model.Teacher, error) {
	md := model.Teacher{}
	err := Db.Raw("select `teacher_id`,`password`,`name`,`college`,`major`,`phone_number` from `t_teacher` where `name` = ? limit 1", username).Scan(&md).Error
	return md, err
}

func (t *TeacherMySQL) QueryAllByID(teacherID string) (model.Teacher, error) {
	md := model.Teacher{}
	err := Db.Raw("select `teacher_id`,`password`,`name`,`college`,`major`,`phone_number` from `t_teacher` where `teacher_id` = ? limit 1", teacherID).Scan(&md).Error
	return md, err
}

func (t *TeacherMySQL) UpdatePassWord(teacherID string, newPwd string) error {
	return Db.Exec("update `t_teacher` set  `password` = ? where `teacher_id` = ?", newPwd, teacherID).Error
}
