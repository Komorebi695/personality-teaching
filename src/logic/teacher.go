package logic

import (
	"encoding/json"
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/dao/redis"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	TeacherRole int8 = 1
	StudentRole int8 = 2
)

type TeacherService struct {
	CTX *gin.Context
}

type teacherFunc interface {
	CheckTeacherPwd(username string, password string) (string, error)
	StoreSession(session model.SessionValue) (string, error)
	CheckTeacherPermission(sessionKey string) (string, error)
	GetTeacherInfo(teacherID string) (model.TeacherInfoResp, error)
}

var _ teacherFunc = &TeacherService{}

func NewTeacherService(c *gin.Context) *TeacherService {
	return &TeacherService{
		CTX: c,
	}
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

func (t *TeacherService)ChangePwd(teacherID string, req model.ChangePwdReq) error {
	teacher, err := mysql.NewTeacherMysql().QueryAllByID(teacherID)
	if err != nil {
		code.CommonResp(t.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return err
	}
	// 解析出明文密码
	clearOldPwd, err := utils.RsaDecrypt(req.OldPassword)
	if err != nil {
		code.CommonResp(t.CTX, http.StatusOK, code.InvalidParam, code.EmptyData)
		return nil
	}
	// 验证旧密码是否正确
	ok, err := utils.CompareHash(teacher.Password, string(clearOldPwd))
	if err != nil {
		code.CommonResp(t.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return nil
	}
	if !ok {
		code.CommonResp(t.CTX, http.StatusOK, code.UnmatchedPassword, code.EmptyData)
		return nil
	}
	// 解析出新密码的明文
	clearNewPwd, err := utils.RsaDecrypt(req.NewPassword)
	if err != nil {
		code.CommonResp(t.CTX, http.StatusOK, code.InvalidParam, code.EmptyData)
		return nil
	}
	// 新密码加密
	pwd, err := utils.Encryption(string(clearNewPwd))
	if err != nil {
		code.CommonResp(t.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return err
	}
	// 存储到数据库
	if err = mysql.NewTeacherMysql().UpdatePassWord(teacherID, pwd); err != nil {
		code.CommonResp(t.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return err
	}
	code.CommonResp(t.CTX, http.StatusOK, code.Success, code.EmptyData)
	return nil
}

func (t *TeacherService) SearchStudent(studentId string) (map[string]float32, error) {
	//得到学生ID，查询t_student_question表，以ID搜索表，按照作业编号分类。
	stu, err := mysql.NewStudentMySQL().QueryStudentInStuQu(studentId)
	//stu中包含知识点编号,题目编号,分数和答案。
	if err != nil {
		return nil, err
	}
	//计算掌握程度
	resp := utils.StuScoreAverage(stu)
	return resp, nil
}
