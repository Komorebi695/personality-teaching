package controller

import (
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/middle"
	"personality-teaching/src/model"

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
	ok, err := teacherService.CheckPassword(req.UserName, req.Password)
	if err != nil {
		logger.L.Error("teacher service QueryAllInfo error :", zap.Error(err))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	if !ok {
		code.CommonResp(c, http.StatusOK, code.WrongPassword, code.EmptyData)
		return
	}
	//  登录成功，生成session并存储至Redis
	sessionKey, err := teacherService.StoreSession(req.UserName)
	if err != nil {
		logger.L.Error("teacher service StoreSession error :", zap.Error(err))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	c.SetCookie(middle.SessionKey, sessionKey, 0, "", "", false, true)
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}
