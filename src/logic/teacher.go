package logic

import (
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/dao/redis"
	"personality-teaching/src/utils"
)

type TeacherService struct{}

type teacherFunc interface {
	CheckPassword(username string, password string) (bool, error)
	StoreSession(username string) (string, error)
}

var _ teacherFunc = &TeacherService{}

func NewTeacherService() *TeacherService {
	return &TeacherService{}
}

func (t *TeacherService) CheckPassword(username string, password string) (bool, error) {
	teacher, err := mysql.NewTeacherMysql().QueryAllInfo(username)
	if err != nil {
		return false, err
	}
	if teacher.TeacherID == "" {
		return false, nil
	}
	// 比较密码
	return utils.CompareHash(teacher.Password, password)
}

// StoreSession 存储session至Redis，返回session_key
func (t *TeacherService) StoreSession(username string) (string, error) {
	teacher, err := mysql.NewTeacherMysql().QueryAllInfo(username)
	if err != nil {
		return "", err
	}
	//  密码匹配将生成session_key并存入Redis
	hashedStr, err := utils.Encryption(teacher.TeacherID)
	if err != nil {
		return "", err
	}
	if err = redis.SetSessionNX(hashedStr, teacher.TeacherID); err != nil {
		return "", err
	}
	return hashedStr, nil
}

// CheckPermission 检查session_key是否具有教师权限
func CheckPermission(sessionKey string) (string, error) {
	return redis.GetTeacherIX(sessionKey)
}
