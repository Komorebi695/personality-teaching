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

type QuestionTypeController struct{}

func QuestionTypeRegister(group *gin.RouterGroup) {
	questionType := &QuestionTypeController{}
	group.GET("/type_query", questionType.TypeQuery)
}

var QuestionTypeService = logic.NewQuestionTypeService()

// TypeQuery godoc
// @Summary 题目类型查找
// @Description 题目类型查找
// @Tags 题目类型管理
// @ID /type/type_query
// @Accept  json
// @Produce  json
// @Param id query int true "题目类型ID"
// @Success 200 {object} code.RespMsg{data=model.TypeQueryOutput} "success"
// @Router /type/type_query [get]
func (point *QuestionTypeController) TypeQuery(c *gin.Context) {
	//从上下文获取参数并校验
	params := &model.TypeQueryInput{}
	if err := c.ShouldBind(params); err != nil {
		logger.L.Error("Input params error:", zap.Error(err))
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	out, err := QuestionTypeService.TypeQueryService(c, params)
	if err == gorm.ErrRecordNotFound {
		code.CommonResp(c, http.StatusInternalServerError, code.RecordNotFound, code.EmptyData)
	} else if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
	}
	code.CommonResp(c, http.StatusOK, code.Success, out)
}
