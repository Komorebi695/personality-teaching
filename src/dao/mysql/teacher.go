package mysql

import "personality-teaching/src/model"

type TeacherMySQL struct {
}

type teacherFunc interface {
	QueryAllInfo(username string) (model.Teacher, error)
}

var _ teacherFunc = &TeacherMySQL{}

func NewTeacherMysql() *TeacherMySQL {
	return &TeacherMySQL{}
}

func (t *TeacherMySQL) QueryAllInfo(username string) (model.Teacher, error) {
	md := model.Teacher{}
	err := db.Raw("select `teacher_id`,`password`,`name`,`college`,`major`,`phone_number` from `t_teacher` where `name` = ? limit 1", username).Scan(&md).Error
	return md, err
}
