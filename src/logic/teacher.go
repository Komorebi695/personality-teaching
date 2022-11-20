package logic

import (
	"encoding/json"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/dao/redis"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

const (
	TeacherRole int8 = 1
	StudentRole int8 = 2
)

type TeacherService struct{}

type teacherFunc interface {
	CheckTeacherPwd(username string, password string) (string, error)
	StoreSession(session model.SessionValue) (string, error)
	CheckTeacherPermission(sessionKey string) (string, error)
	GetTeacherInfo(teacherID string) (model.TeacherInfoResp, error)
}

var _ teacherFunc = &TeacherService{}

func NewTeacherService() *TeacherService {
	return &TeacherService{}
}

// CheckTeacherPwd  校验通过返回teacherID，失败返回空字符串
func (t *TeacherService) CheckTeacherPwd(username string, password string) (string, error) {
	teacher, err := mysql.NewTeacherMysql().QueryAllByName(username)
	if err != nil || teacher.TeacherID == "" {
		return "", err
	}
	// 比较密码
	ok, err := utils.CompareHash(teacher.Password, password)
	if err != nil || !ok {
		return "", err
	}
	return teacher.TeacherID, nil
}

// StoreSession 存储session至Redis，返回session_key
func (t *TeacherService) StoreSession(session model.SessionValue) (string, error) {
	sessionKey := utils.GetUUID()
	byteData, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	if err = redis.SetSessionNX(sessionKey, string(byteData)); err != nil {
		return "", err
	}
	return sessionKey, nil
}

// CheckTeacherPermission 检查session_key是否具有教师权限，成功返回teacher_id并重新设置过期时间
func (t *TeacherService) CheckTeacherPermission(sessionKey string) (string, error) {
	sessionValue, err := redis.GetSessionValue(sessionKey)
	if err != nil {
		return "", err
	}
	if sessionValue == "" {
		return "", err
	}
	//  反序列化至结构体
	var sv model.SessionValue
	if err = json.Unmarshal([]byte(sessionValue), &sv); err != nil {
		return "", err
	}
	if sv.RoleType != TeacherRole {
		return "", err
	}
	//  重置键过期时间
	if err = redis.ResetExpireTime(sessionKey); err != nil {
		return "", err
	}
	return sv.UserID, nil
}

func (*TeacherService) GetTeacherInfo(teacherID string) (model.TeacherInfoResp, error) {
	t, err := mysql.NewTeacherMysql().QueryAllByID(teacherID)
	if err != nil {
		return model.TeacherInfoResp{}, err
	}
	return model.TeacherInfoResp{
		TeacherID:   t.TeacherID,
		Name:        t.Name,
		College:     t.College,
		Major:       t.Major,
		PhoneNumber: t.PhoneNumber,
	}, nil
}
