package controller

import (
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func AddClass(c *gin.Context) {
	var req model.ClassAddReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	teacherID := c.GetString(utils.TeacherID)
	class, err := logic.NewClassService().ClassAdd(teacherID, req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("add class error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, class)
}

func UpdateClass(c *gin.Context) {
	var req model.ClassUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	//  校验教师是否有修改此班级权限
	teacherID := c.GetString(utils.TeacherID)
	legal, err := logic.NewClassService().CheckPermission(teacherID, req.ClassID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("check teacher's permission with class error: ", zap.Error(err))
		return
	}
	if !legal {
		code.CommonResp(c, http.StatusOK, code.InvalidPermission, code.EmptyData)
		return
	}
	//  执行更新
	if err := logic.NewClassService().ClassUpdate(req); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("update class error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

func DeleteClass(c *gin.Context) {
	var req model.ClassDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	teacherID := c.GetString(utils.TeacherID)
	//  校验教师是否有修改此班级权限
	legal, err := logic.NewClassService().CheckPermission(teacherID, req.ClassID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("check teacher's permission with class error: ", zap.Error(err))
		return
	}
	if !legal {
		code.CommonResp(c, http.StatusOK, code.InvalidPermission, code.EmptyData)
		return
	}
	// 执行删除
	if err := logic.NewClassService().ClassDelete(teacherID, req.ClassID); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("delete class error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

func ClassInfo(c *gin.Context) {
	var req model.ClassInfoReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	class, err := logic.NewClassService().ClassInfo(req.ClassID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("get class info error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, class)
}

func ClassList(c *gin.Context) {
	var req model.ClassListReq
	if err := c.ShouldBind(&req); err != nil {
		code.RespList(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData, 0)
		return
	}
	teacherID := c.GetString(utils.TeacherID)
	classes, total, err := logic.NewClassService().ClassInfoList(teacherID, req)
	if err != nil {
		code.RespList(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData, 0)
		logger.L.Error("get class list error: ", zap.Error(err))
		return
	}

	code.RespList(c, http.StatusOK, code.Success, classes, total)
}

// ClassNameCheck 检查班级名是否存在
func ClassNameCheck(c *gin.Context) {
	var req model.ClassNameReq
	if err := c.ShouldBind(&req); err != nil {
		code.RespList(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData, 0)
		return
	}
	flag, err := logic.NewClassService().ClassNameCheck(req.Name)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("check class name error: ", zap.Error(err))
		return
	} else if flag {
		code.CommonResp(c, http.StatusOK, code.ClassNameExit, code.EmptyData)
		return
	}
	code.CommonResp(c, http.StatusOK, code.RecordNotFound, code.EmptyData)
	return
}
