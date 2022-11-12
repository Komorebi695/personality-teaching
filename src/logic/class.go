package logic

import (
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

type ClassService struct{} //所有size为0的变量都用的是同一块内存  zerobase

func NewClassService() *ClassService {
	return &ClassService{}
}

func (c *ClassService) ClassAdd(teacherID string, req model.ClassAddReq) (model.Class, error) {
	class := model.Class{
		ClassID:       utils.GenSnowID(),
		BaseClassInfo: req.BaseClassInfo,
	}
	if err := mysql.NewClassMysql().InsertClass(teacherID, class); err != nil {
		return model.Class{}, err
	}
	return class, nil
}

func (c *ClassService) ClassUpdate(req model.ClassUpdateReq) error {
	class := model.Class{
		ClassID:       req.ClassID,
		BaseClassInfo: req.BaseClassInfo,
	}
	return mysql.NewClassMysql().UpdateClass(class)
}

func (c *ClassService) ClassDelete(teacherID string, classID string) error {
	return mysql.NewClassMysql().DeleteClass(teacherID, classID)
}

func (c *ClassService) CheckPermission(teacherID, classID string) (bool, error) {
	return mysql.NewClassMysql().CheckTeacherClass(teacherID, classID)
}

func (c *ClassService) ClassInfoList(teacherID string, req model.ClassListReq) ([]model.Class, error) {
	return mysql.NewClassMysql().QueryClassList(teacherID, req)
}

func (c *ClassService) ClassInfo(classID string) (model.Class, error) {
	return mysql.NewClassMysql().QueryClass(classID)
}
