package controller

import (
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SearchExam ,搜索试卷
// Param:
//
//	text: 搜索文本
//
// Router /teacher/exam/search [post]
func SearchExam(c *gin.Context) {
	var req model.SearchReq
	// 绑定请求参数到 ExamAddReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		logger.L.Error("add exam error: ", zap.Error(err))
		return
	}
	// 获取当前登录的老师编号
	teacherID := c.GetString(utils.TeacherID)
	res, err := logic.NewExamService().SearchExam(req.Text, teacherID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("search exam error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, res)
}

// AddExam ,新增试卷
// Param:
//
//	exam_name: 试卷名称
//	questions: 试题
//	comment: 备注
//
// Router /teacher/exam [post]
func AddExam(c *gin.Context) {
	var req model.ExamAddReq
	// 绑定请求参数到 ExamAddReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		logger.L.Error("add exam error: ", zap.Error(err))
		return
	}
	// 获取当前登录的老师编号
	teacherID := c.GetString(utils.TeacherID)
	exam, err := logic.NewExamService().Add(teacherID, req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("add exam error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, exam)
}

// UpdateExam ,更新试卷内容
// Param:
//
//	exam_id: 试卷编号
//	exam_name: 试卷名称
//	questions: 试题
//	comment: 备注
//
// Router /teacher/exam [put]
func UpdateExam(c *gin.Context) {
	var req model.ExamUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	// 执行更新
	if err := logic.NewExamService().Update(req); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("update exam error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

// DeleteExam ,删除试卷
// Param:
//
//	exam_id: 试卷编号
//
// Router /teacher/exam [delete]
func DeleteExam(c *gin.Context) {
	var req model.ExamDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	// 执行删除
	if err := logic.NewExamService().Delete(req); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("delete exam error：", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

// ExamList ,获取当前登录老师的试卷列表
// Param:
// offset 第几页
// page_size 页面大小
// Router /teacher/exam/list [get]
func ExamList(c *gin.Context) {
	var req model.PagingReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	// 获取当前登录的老师编号
	teacherID := c.GetString(utils.TeacherID)
	// 查询
	resp, err := logic.NewExamService().List(teacherID, req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query exam list error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}

// ExamInfo ,获取试卷详细信息
// Param:
// exam_id 试卷编号
// Router /teacher/exam [get]
func ExamInfo(c *gin.Context) {
	var req model.ExamIDReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	// 查询
	examDetail, err := logic.NewExamService().Details(req.ExamID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query exam detail error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, examDetail)
}

// SendExam ,发布试卷给个人
// Param:
// exam_id 试卷编号
// student_id 学生编号 || class_id 班级编号
// start_time 开始时间
// end_time 结束时间
// Router /teacher/exam/send/1 [post]  -- 个人
// Router /teacher/exam/send/2 [post]  -- 班级
func SendExam(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	// 按个人发放
	if num == 1 {
		var req model.SendPersonReq
		if err := c.ShouldBindJSON(&req); err != nil {
			code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
			logger.L.Error("绑定参数错误: ", zap.Error(err))
			return
		}
		if err := logic.NewExamService().SendPerson(req); err != nil {
			code.CommonResp(c, http.StatusOK, code.ServerBusy, code.EmptyData)
			logger.L.Error("send exam by person error: ", zap.Error(err))
			return
		}
		code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
	} else if num == 2 {
		// 按班级发放
		var req model.SendClassReq
		if err := c.ShouldBind(&req); err != nil {
			code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
			return
		}
		if err := logic.NewExamService().SendClass(req); err != nil {
			code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
			logger.L.Error("send exam by class error: ", zap.Error(err))
			return
		}
		code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
	}
}

// ReleaseStudentList ,获取发布试卷班级的学生
// Param:
// exam_id 试卷编号
// class_id 班级编号
// Router /teacher/exam/student/list
func ReleaseStudentList(c *gin.Context) {
	var req model.ReleaseExamReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	studentList, err := logic.NewExamService().ReleaseStudentList(req.ClassID.ClassID, req.ExamIDReq.ExamID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query release student list error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, studentList)
	return
}

// GetTeacherExamList .获取发布试卷信息
func GetTeacherExamList(c *gin.Context) {
	var req model.GetTeacherExamListReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	ExamId, err := mysql.GetExamIDByStudentID(req.StudentID.StudentID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query exam_id error: ", zap.Error(err))
		return
	}
	examDetail, err := logic.NewExamService().Details(ExamId)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query exam detail error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, examDetail)
}
