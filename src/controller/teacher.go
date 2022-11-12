package controller

import (
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/middle"
	"personality-teaching/src/model"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func TeacherLogin(c *gin.Context) {
	req := model.TeacherLoginReq{}
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	teacherService := logic.NewTeacherService()
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
	c.SetCookie(middle.SessionKey, sessionKey, 0, "", "", false, false)
	code.CommonResp(c, http.StatusOK, code.Success, teacherID)
}
