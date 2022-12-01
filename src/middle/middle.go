package middle

import (
	"fmt"
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// VerifyTeacher 负责验证用户是否有教师权限，若有则在上下文中存入teacher_id
func VerifyTeacher(c *gin.Context) {
	key, err := c.Cookie(utils.SessionKey)
	if err == http.ErrNoCookie {
		code.CommonResp(c, http.StatusOK, code.NeedLogin, code.EmptyData)
		return
	}
	teacherID, err := logic.NewTeacherService(c).CheckTeacherPermission(key)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("check teacher permission error: ", zap.Error(err), zap.String("session", key))
		return
	}
	if teacherID == "" {
		code.CommonResp(c, http.StatusOK, code.InvalidPermission, code.EmptyData)
		return
	}
	c.Set(utils.TeacherID, teacherID)
}

func VerifyStudent(c *gin.Context) {
	key, err := c.Cookie(utils.SessionKey)
	if err == http.ErrNoCookie {
		code.CommonResp(c, http.StatusOK, code.NeedLogin, code.EmptyData)
		return
	}
	studentID, err := logic.NewStudentService(c).CheckStudentPermission(key)
	if err != nil {
		logger.L.Error("check student permission error: ", zap.Error(err), zap.String("session", key))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	if studentID == "" {
		code.CommonResp(c, http.StatusOK, code.InvalidPermission, code.EmptyData)
		return
	}
	c.Set(utils.StudentID, studentID)
}

// VerifyAny 验证是老师或者学生期中一者即可
func VerifyAny(c *gin.Context) {
	key, err := c.Cookie(utils.SessionKey)
	if err == http.ErrNoCookie {
		code.CommonResp(c, http.StatusOK, code.NeedLogin, code.EmptyData)
		return
	}
	var teacherID, studentID string
	teacherID, err = logic.NewTeacherService(c).CheckTeacherPermission(key)
	if err != nil {
		logger.L.Error("check student and teacher permission error: ", zap.Error(err), zap.String("session", key))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	studentID, err = logic.NewStudentService(c).CheckStudentPermission(key)
	if err != nil {
		logger.L.Error("check student and teacher permission error: ", zap.Error(err), zap.String("session", key))
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		return
	}
	if studentID == "" && teacherID != "" {
		c.Set(utils.Role, utils.TeacherID)
		c.Set(utils.TeacherID, teacherID)
		return
	} else if studentID != "" && teacherID == "" {
		c.Set(utils.Role, utils.StudentID)
		c.Set(utils.StudentID, studentID)
		return
	} else {
		code.CommonResp(c, http.StatusOK, code.NeedLogin, code.EmptyData)
	}
}

func ChangePassword(c *gin.Context) {
	var req model.ChangePwdReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err)
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	var err error
	if c.GetString(utils.Role) == utils.TeacherID {
		err = logic.NewTeacherService(c).ChangePwd(c.GetString(utils.TeacherID), req)
	} else if c.GetString(utils.Role) == utils.StudentID {
		err = logic.NewStudentService(c).ChangePwd(c.GetString(utils.StudentID), req)
	}
	if err != nil {
		logger.L.Error("change password error: ", zap.Error(err))
	}
}
