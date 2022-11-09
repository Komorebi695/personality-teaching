package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/model"
)

type QuestionController struct{}

func QuestionRegister(group *gin.RouterGroup) {
	question := &QuestionController{}
	group.GET("/question_list", question.QuestionList)
	group.DELETE("/question_delete", question.QuestionDelete)
	group.GET("/question_detail", question.QuestionDetail)
	group.POST("/question_add", question.QuestionAdd)
	group.PUT("/question_update", question.QuestionUpdate)

}

var questionService = logic.NewQuestionService()

// QuestionList godoc
// @Summary 题目列表
// @Description 题目列表
// @Tags 题目管理
// @ID /question/question_list
// @Accept  json
// @Produce  json
// @Param context query string false "题目关键词"
// @Param type query int true "题目类型"
// @Param level query int true "困难程度"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} code.RespMsg{data=model.QuestionListOutput} "success"
// @Router /question/question_list [get]
func (question *QuestionController) QuestionList(c *gin.Context) {
	//从上下文获取参数并校验
	params := &model.QuestionListInput{}
	if err := c.ShouldBind(params); err != nil {
		logger.L.Error("Input params error:", zap.Error(err))
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	out, err := questionService.QuestionListService(c, params)
	if err == gorm.ErrRecordNotFound {
		code.CommonResp(c, http.StatusInternalServerError, code.RecordNotFound, code.EmptyData)
	} else if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
	}
	code.CommonResp(c, http.StatusOK, code.Success, out)
}

// QuestionDelete godoc
// @Summary 题目删除
// @Description 题目删除
// @Tags 题目管理
// @ID /question/question_delete
// @Accept  json
// @Produce  json
// @Param question_id query string true "题目ID"
// @Success 200 {object} code.RespMsg{data=string} "success"
// @Router /question/question_delete [delete]
func (question *QuestionController) QuestionDelete(c *gin.Context) {
	params := &model.QuestionDeleteInput{}
	if err := c.ShouldBind(params); err != nil {
		logger.L.Error("Input params error:", zap.Error(err))
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	err := questionService.QuestionDeleteService(c, params)
	if err == gorm.ErrRecordNotFound {
		code.CommonResp(c, http.StatusInternalServerError, code.RecordNotFound, code.EmptyData)
	} else if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

// QuestionAdd godoc
// @Summary 添加题目
// @Description 添加题目
// @Tags 题目管理
// @ID /question/question_add
// @Accept  json
// @Produce  json
// @Param body body model.QuestionAddInput true "body"
// @Success 200 {object} code.RespMsg{data=string} "success"
// @Router /question/question_add [post]
func (question *QuestionController) QuestionAdd(c *gin.Context) {
	params := &model.QuestionAddInput{}
	if err := c.ShouldBind(params); err != nil {
		logger.L.Error("Input params error:", zap.Error(err))
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	err := questionService.QuestionAddService(c, params)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}

// QuestionDetail godoc
// @Summary 题目详情
// @Description 题目详情
// @Tags 题目管理
// @ID /question/question_detail
// @Accept  json
// @Produce  json
// @Param question_id query string true "题目ID"
// @Success 200 {object} code.RespMsg{data=mysql.QuestionDetail} "success"
// @Router /question/question_detail [get]
func (question *QuestionController) QuestionDetail(c *gin.Context) {
	params := &model.QuestionDetailInput{}
	if err := c.ShouldBind(params); err != nil {
		logger.L.Error("Input params error:", zap.Error(err))
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	questionDetail, err := questionService.QuestionDetailService(c, params)
	if err == gorm.ErrRecordNotFound {
		code.CommonResp(c, http.StatusInternalServerError, code.RecordNotFound, code.EmptyData)
	} else if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
	}
	code.CommonResp(c, http.StatusOK, code.Success, questionDetail)
}

// QuestionUpdate godoc
// @Summary 修改题目
// @Description 修改题目
// @Tags 题目管理
// @ID /question/question_update
// @Accept  json
// @Produce  json
// @Param body body model.QuestionUpdateInput true "body"
// @Success 200 {object} code.RespMsg{data=string} "success"
// @Router /question/question_update [put]
func (question *QuestionController) QuestionUpdate(c *gin.Context) {
	params := &model.QuestionUpdateInput{}
	if err := c.ShouldBind(params); err != nil {
		logger.L.Error("Input params error:", zap.Error(err))
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	err := questionService.QuestionUpdateService(c, params)
	if err == gorm.ErrRecordNotFound {
		code.CommonResp(c, http.StatusInternalServerError, code.RecordNotFound, code.EmptyData)
	} else if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
	}
	code.CommonResp(c, http.StatusOK, code.Success, code.EmptyData)
}
