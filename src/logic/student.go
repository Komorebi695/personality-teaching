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

type studentFunc interface {
	CreateStudent(req model.CreateStudentReq) (model.CreateStudentResp, error)
	UpdateClassID(req model.AddStudentToClassReq) (model.AddStudentClassResp, error)
	GetStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, int, error)
	RemoveStudentClass(studentID string) error
	CheckStudentClass(studentID string, classID string) (bool, error)
	CheckStudentPermission(sessionKey string) (string, error)
	CheckPwd(req model.LoginReq) (string, error)
}

var _ studentFunc = &StudentService{}

type StudentService struct {
	CTX *gin.Context
}

func NewStudentService(c *gin.Context) *StudentService {
	return &StudentService{
		CTX: c,
	}
}

func (s *StudentService) CreateStudent(req model.CreateStudentReq) (model.CreateStudentResp, error) {
	student := model.Student{
		BaseStudentInfo: req.BaseStudentInfo,
		StudentID:       utils.GenSnowID(),
		Password:        utils.GetDefaultPassWord(),
		ClassID:         utils.EmptyClassID,
	}
	if err := mysql.NewStudentMySQL().InsertStudent(student); err != nil {
		return model.CreateStudentResp{}, err
	}
	return model.CreateStudentResp{
		StudentID:       student.StudentID,
		BaseStudentInfo: student.BaseStudentInfo,
	}, nil
}

func (s *StudentService) UpdateClassID(req model.AddStudentToClassReq) (model.AddStudentClassResp, error) {
	if err := mysql.NewStudentMySQL().UpdateClassID(req.StudentID, req.ClassID); err != nil {
		return model.AddStudentClassResp{}, err
	}
	student, err := mysql.NewStudentMySQL().QueryStudent(req.StudentID)
	if err != nil {
		return model.AddStudentClassResp{}, err
	}
	return model.AddStudentClassResp{
		StudentID:       student.StudentID,
		ClassID:         student.ClassID,
		BaseStudentInfo: student.BaseStudentInfo,
	}, nil
}

func (s *StudentService) GetStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, int, error) {
	students, total, err := mysql.NewStudentMySQL().QueryStudentsInClass(req)
	if err != nil {
		return []model.ClassStudentListResp{}, 0, err
	}
	return students, total, nil
}

func (s *StudentService) RemoveStudentClass(studentID string) error {
	return mysql.NewStudentMySQL().UpdateClassID(studentID, utils.EmptyClassID)
}

func (s *StudentService) CheckStudentClass(studentID string, classID string) (bool, error) {
	return mysql.NewStudentMySQL().CheckStudentClass(studentID, classID)
}

// CheckStudentPermission ??????????????????studentID
func (s *StudentService) CheckStudentPermission(sessionKey string) (string, error) {
	sessionVal, err := redis.GetSessionValue(sessionKey)
	if err != nil {
		return "", err
	}
	if sessionVal == "" {
		return "", nil
	}
	var sv model.SessionValue
	if err = json.Unmarshal([]byte(sessionVal), &sv); err != nil {
		return "", err
	}
	if sv.RoleType != StudentRole {
		return "", err
	}
	//  ?????????????????????
	if err = redis.ResetExpireTime(sessionKey); err != nil {
		return "", err
	}
	return sv.UserID, nil
}

func (s *StudentService) ChangePwd(studentID string, req model.ChangePwdReq) error {
	student, err := mysql.NewStudentMySQL().QueryStudent(studentID)
	if err != nil {
		code.CommonResp(s.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return err
	}
	// ?????????????????????
	clearOldPwd, err := utils.RsaDecrypt(req.OldPassword)
	if err != nil {
		code.CommonResp(s.CTX, http.StatusOK, code.InvalidParam, code.EmptyData)
		return nil
	}
	// ???????????????????????????
	ok, err := utils.CompareHash(student.Password, string(clearOldPwd))
	if err != nil {
		code.CommonResp(s.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return nil
	}
	if !ok {
		code.CommonResp(s.CTX, http.StatusOK, code.UnmatchedPassword, code.EmptyData)
		return nil
	}
	// ???????????????????????????
	clearNewPwd, err := utils.RsaDecrypt(req.NewPassword)
	if err != nil {
		code.CommonResp(s.CTX, http.StatusOK, code.InvalidParam, code.EmptyData)
		return nil
	}
	// ???????????????
	pwd, err := utils.Encryption(string(clearNewPwd))
	if err != nil {
		code.CommonResp(s.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return err
	}
	// ??????????????????
	if err = mysql.NewStudentMySQL().UpdatePassWord(studentID, pwd); err != nil {
		code.CommonResp(s.CTX, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return err
	}
	code.CommonResp(s.CTX, http.StatusOK, code.Success, code.EmptyData)
	return nil
}

func (s *StudentService) CheckPwd(req model.LoginReq) (string, error) {
	student, err := mysql.NewStudentMySQL().QueryAllByName(req.UserName)
	if err != nil || student.StudentID == "" {
		return "", err
	}
	ok, err := utils.CompareHash(student.Password, req.Password)
	if err != nil || !ok {
		return "", err
	}
	return student.StudentID, nil
}
