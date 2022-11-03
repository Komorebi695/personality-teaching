package logic

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
)

type QuestionTypeService struct{}

type questionTypeFunc interface {
	TypeQueryService(c *gin.Context, params *model.TypeQueryInput) (*model.TypeQueryOutput, error)
}

var _ questionTypeFunc = &QuestionTypeService{}

func NewQuestionTypeService() *QuestionTypeService {
	return &QuestionTypeService{}
}

func (q *QuestionTypeService) TypeQueryService(c *gin.Context, params *model.TypeQueryInput) (*model.TypeQueryOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionTypeService` -> get pool err:", zap.Error(err))
		return nil, err
	}
	questionType := &mysql.TQuestionType{Id: params.Id}
	find, err := questionType.Find(c, tx)
	if err != nil {
		logger.L.Error("`QuestionTypeService` -> questionType.Find err:", zap.Error(err))
		return nil, err
	}
	out := &model.TypeQueryOutput{
		Id:   find.Id,
		Name: find.Name,
	}
	return out, nil
}
