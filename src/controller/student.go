package controller

import (
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreateStudent(c *gin.Context) {
	var req model.CreateStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}

	resp, err := logic.NewStudentService(c).CreateStudent(req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("CreateStudent error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}

func AddStudentToClass(c *gin.Context) {
	var req model.AddStudentToClassReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	// 校验教师是否有该班级的权限和班级是否存在
	teacherID := c.GetString(utils.TeacherID)
	ok, err := logic.NewClassService().CheckPermission(teacherID, req.ClassID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("CheckPermission error: ", zap.Error(err))
		return
	}
	if !ok {
		code.CommonResp(c, http.StatusOK, code.InvalidPermission, code.EmptyData)
		return
	}
	// 将学生加入班级
	resp, err := logic.NewStudentService(c).UpdateClassID(req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("UpdateClassID error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}

func StudentsInClass(c *gin.Context) {
	var req model.ClassStudentListReq
	if err := c.ShouldBind(&req); err != nil {
		code.RespList(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData, 0)
		return
	}
	// 校验教师是否有该班级的权限和班级是否存在
	teacherID := c.GetString(utils.TeacherID)
	ok, err := logic.NewClassService().CheckPermission(teacherID, req.ClassID)
	if err != nil {
		code.RespList(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData, 0)
		logger.L.Error("CheckPermission error: ", zap.Error(err))
		return
	}
	if !ok {
		code.RespList(c, http.StatusOK, code.InvalidPermission, code.EmptyData, 0)
		return
	}
	resp, total, err := logic.NewStudentService(c).GetStudentsInClass(req)
	if err != nil {
		code.RespList(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData, 0)
		logger.L.Error("GetStudentInClass error: ", zap.Error(err))
		return
	}
	code.RespList(c, http.StatusOK, code.Success, resp, total)
}

// StudentNotInClass 查询未加入班级的学生
func StudentNotInClass(c *gin.Context) {
	var req model.EmptyClassStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.RespList(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData, 0)
		return
	}
	resp, total, err := logic.NewStudentService(c).GetStudentsNotInClass(model.ClassStudentListReq{
		ClassID:  utils.EmptyClassID,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}, req.Content)
	if err != nil {
		code.RespList(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData, 0)
		logger.L.Error("GetStudentInClass error: ", zap.Error(err))
		return
	}
	code.RespList(c, http.StatusOK, code.Success, resp, total)
}

func DeleteClassStudent(c *gin.Context) {
	var req model.DeleteClassStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	// 校验班级学生关系
	ok, err := logic.NewStudentService(c).CheckStudentClass(req.StudentID, req.ClassID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("CheckStudentClass error: ", zap.Error(err))
		return
	}
	if !ok {
		code.CommonResp(c, http.StatusOK, code.NotInClass, code.EmptyData)
		return
	}
	if err = logic.NewStudentService(c).RemoveStudentClass(req.StudentID); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("RemoveStudentClass error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

func StudentLogin(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	// 解析密码明文
	plaintext, err := utils.RsaDecrypt(req.Password)
	if err != nil {
		code.CommonResp(c, http.StatusOK, code.WrongPassword, code.EmptyData)
		return
	}
	req.Password = string(plaintext)

	studentService := logic.NewStudentService(c)
	studentID, err := studentService.CheckPwd(req)
	if err != nil {
		logger.L.Error("student login error :", zap.Error(err))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	if studentID == "" {
		code.CommonResp(c, http.StatusOK, code.WrongPassword, code.EmptyData)
		return
	}
	//  登录成功，生成session并存储至Redis
	session := model.SessionValue{
		UserID:     studentID,
		RoleType:   logic.StudentRole,
		CreateTime: time.Now().Unix(),
	}
	sessionKey, err := logic.NewTeacherService(c).StoreSession(session)
	if err != nil {
		logger.L.Error("teacher service StoreSession error :", zap.Error(err))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	c.SetCookie(utils.SessionKey, sessionKey, 0, "", "", false, false)
	code.CommonResp(c, http.StatusOK, code.Success, studentID)
}

func SearchStudent(c *gin.Context) {
	var req model.SearchStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	studentList, err := logic.NewStudentService(c).SearchStudent(req.SearchText)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("search student error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, studentList)
	return
}

func DeleteStudent(c *gin.Context) {
	var req model.DeleteStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}

	if err := logic.NewStudentService(c).RemoveStudent(req.StudentID); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("RemoveStudent error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

func UpdateStudent(c *gin.Context) {
	var req model.CreateStudentResp
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}

	if err := logic.NewStudentService(c).UpdateStudent(req); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("Update student info error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}
