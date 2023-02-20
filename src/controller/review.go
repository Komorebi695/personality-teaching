package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/model"
)

// ReviewUpdate ,更新评阅的试卷
// Router /teacher/review [put]
func ReviewUpdate(c *gin.Context) {
	var req model.ReviewUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	if err := logic.NewReviewService().UpdateReview(req); err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("update review error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

// ReviewClass ,查询做该试卷的班级
//Param：
// exam_id：试卷编号
// Router /teacher/review/class [get]
func ReviewClass(c *gin.Context) {
	var req model.ExamIDReq
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err.Error())
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}

	classList, err := logic.NewReviewService().QueryReviewClass(req.ExamID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query exam class list error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, classList)
}

// ReviewStudentList ,查询做该试卷对应班级的学生
//Param：
// exam_id：试卷编号
// class_id：班级编号
// Router /teacher/review/student/list [get]
func ReviewStudentList(c *gin.Context) {
	var req model.ReviewStudentListReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	studentList, err := logic.NewReviewService().QueryReviewStudent(req.ClassID, req.ExamID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query exam student list error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, studentList)
}

// ReviewStudent ,查询试卷答案
// exam_id：试卷编号
// student_id：学生编号
// Router /teacher/review/student [get]
func ReviewStudent(c *gin.Context) {
	var req model.ReviewStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusOK, code.InvalidParam, code.EmptyData)
		return
	}
	studentAnswer, err := logic.NewReviewService().QueryStudentAnswer(req.ExamID, req.StudentID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("query student exam answer error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, studentAnswer)

}
