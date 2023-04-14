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

func TeacherLogin(c *gin.Context) {
	req := model.LoginReq{}
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	// 解析密码明文
	plaintext, err := utils.RsaDecrypt(req.Password)
	if err != nil {
		logger.L.Error("RsaDecrypt error :", zap.Error(err))
		code.CommonResp(c, http.StatusOK, code.WrongPassword, code.EmptyData)
		return
	}
	req.Password = string(plaintext)

	teacherService := logic.NewTeacherService(c)
	teacherID, err := teacherService.CheckTeacherPwd(req.UserName, req.Password)
	if err != nil {
		logger.L.Error("teacher service QueryAllInfo error :", zap.Error(err))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	if teacherID == "" {
		code.CommonResp(c, http.StatusOK, code.WrongPassword, code.EmptyData)
		return
	}
	//  登录成功，生成session并存储至Redis
	session := model.SessionValue{
		UserID:     teacherID,
		RoleType:   logic.TeacherRole,
		CreateTime: time.Now().Unix(),
	}
	sessionKey, err := teacherService.StoreSession(session)
	if err != nil {
		logger.L.Error("teacher service StoreSession error :", zap.Error(err))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	c.SetCookie(utils.SessionKey, sessionKey, 0, "", "", false, false)
	code.CommonResp(c, http.StatusOK, code.Success, teacherID)
}

func TeacherInfo(c *gin.Context) {
	teacherID := c.GetString(utils.TeacherID)

	resp, err := logic.NewTeacherService(c).GetTeacherInfo(teacherID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("GetTeacherInfo error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}

// 根据学号查询学生作业掌握情况。
func TeacherSearchStudentID(c *gin.Context) {
	studentresp := model.SearchStudentIDResp{}
	if err := c.ShouldBind(&studentresp); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	studentID := studentresp.StudentID
	resp, err := logic.NewTeacherService(c).SearchStudent(studentID)
	if err != nil {
		logger.L.Error("TeacherSearchStudentID error: ", zap.Error(err))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}
