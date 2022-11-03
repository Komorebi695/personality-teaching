package logic

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
	"time"
)

type QuestionService struct{}

type questionFunc interface {
	QuestionListService(c *gin.Context, params *model.QuestionListInput) (*model.QuestionListOutput, error)
	QuestionDeleteService(c *gin.Context, params *model.QuestionDeleteInput) error
	QuestionAddService(c *gin.Context, params *model.QuestionAddInput) error
	QuestionDetailService(c *gin.Context, params *model.QuestionDetailInput) (*mysql.QuestionDetail, error)
}

var _ questionFunc = &QuestionService{}

func NewQuestionService() *QuestionService {
	return &QuestionService{}
}

func (q *QuestionService) QuestionListService(c *gin.Context, params *model.QuestionListInput) (*model.QuestionListOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionListService` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//从db中分页读取基本信息
	questionInfo := &mysql.TQuestion{}
	list, total, err := questionInfo.PageList(c, tx, params)
	if err != nil {
		logger.L.Error("`QuestionListService` -> questionInfo.PageList err:", zap.Error(err))
		return nil, err
	}
	//格式化输出信息
	var outList []model.QuestionListItemOutput
	for _, listItem := range list {
		outItem := model.QuestionListItemOutput{
			QuestionId:   listItem.QuestionId,
			QuestionName: listItem.Name,
			Context:      listItem.Context,
			Answer:       listItem.Answer,
			Type:         listItem.Type,
			Level:        listItem.Level,
			CreateUser:   listItem.CreateUser,
		}
		outList = append(outList, outItem)
	}
	out := &model.QuestionListOutput{
		Total: total,
		List:  outList,
	}
	return out, nil
}

func (q *QuestionService) QuestionDeleteService(c *gin.Context, params *model.QuestionDeleteInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionDeleteService` -> get pool err:", zap.Error(err))
		return err
	}
	//读取基本信息
	questionInfo := &mysql.TQuestion{QuestionId: params.QuestionId}
	questionInfo, err = questionInfo.Find(c, tx)
	if err != nil {
		logger.L.Error("`QuestionDeleteService` -> TQuestion.FindById err:", zap.Error(err))
		return err
	}
	questionInfo.IsDelete = 1
	if err = questionInfo.Save(c, tx); err != nil {
		logger.L.Error("`QuestionDeleteService` -> TQuestion.Save err:", zap.Error(err))
		return err
	}
	return nil
}

func (q *QuestionService) QuestionAddService(c *gin.Context, params *model.QuestionAddInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionAddService` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//判断题目是否重复插入
	questionInfo := &mysql.TQuestion{Name: params.QuestionName}
	if _, err = questionInfo.Find(c, tx); err == nil {
		tx.Rollback()
		logger.L.Error("`QuestionAddService` -> The problem's name already exists:", zap.Error(err))
		return err
	}
	//包装题目信息
	questionModel := &mysql.TQuestion{
		//使用雪花ID生成
		QuestionId: utils.GenSnowID(),
		Name:       params.QuestionName,
		Level:      params.Level,
		Type:       params.Type,
		Context:    params.Context,
		Answer:     params.Answer,
		CreateUser: params.CreateUser,
		IsDelete:   0,
	}
	if err = questionModel.Save(c, tx); err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionAddService` -> questionModel.Save err:", zap.Error(err))
		return err
	}

	//循环依次插入知识点题目关联表
	knpIdList := params.GetKnpIdByModel()
	for _, s := range knpIdList {
		knowledgePointQuestion := &mysql.TKnowledgePointQuestion{
			KnpId:      s,
			QuestionId: questionModel.QuestionId,
		}
		if err = knowledgePointQuestion.Save(c, tx); err != nil {
			tx.Rollback()
			logger.L.Error("`QuestionAddService` -> knowledgePointQuestion.Save err:", zap.Error(err))
			return err
		}
	}
	//若是选择题，插入选项表内容
	if params.Type == 1 {
		for _, option := range params.QuestionOptionList {
			questionOption := &mysql.TQuestionOption{
				QuestionId: questionModel.QuestionId,
				Context:    option.Context,
				IsAnswer:   option.IsAnswer,
			}
			if err = questionOption.Save(c, tx); err != nil {
				tx.Rollback()
				logger.L.Error("`QuestionAddService` -> TQuestionOption.Save err:", zap.Error(err))
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (q *QuestionService) QuestionDetailService(c *gin.Context, params *model.QuestionDetailInput) (*mysql.QuestionDetail, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionDetailService` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//获取问题详情
	questionInfo := &mysql.TQuestion{QuestionId: params.QuestionId}
	questionInfo, err = questionInfo.Find(c, tx)
	if err != nil {
		logger.L.Error("`QuestionDetailService` -> questionInfo.FindById err:", zap.Error(err))
		return nil, err
	}
	//题目对应的知识点编号表
	questionPointSearch := &mysql.TKnowledgePointQuestion{QuestionId: questionInfo.QuestionId}
	questionPointList, err := questionPointSearch.Find(c, tx, questionPointSearch)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.L.Error("`QuestionDetailService` -> TKnowledgePointQuestion.FindById err:", zap.Error(err))
		return nil, err
	}
	//根据编号表查询知识点列表
	var knowledgePointList []*mysql.TKnowledgePoint
	for _, point := range questionPointList {
		pointSearch := &mysql.TKnowledgePoint{KnpId: point.KnpId}
		pointSearch, err = pointSearch.FindById(c, tx)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.L.Error("`QuestionDetailService` -> TKnowledgePoint.FindById err:", zap.Error(err))
			return nil, err
		}
		knowledgePointList = append(knowledgePointList, pointSearch)
	}
	//若是选择题，查询选项表
	var optionList []*mysql.TQuestionOption
	if questionInfo.Type == 1 {
		optionItem := &mysql.TQuestionOption{QuestionId: params.QuestionId}
		optionList, err = optionItem.FindByQuestionId(c, tx)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.L.Error("`QuestionDetailService` -> TQuestionOption.FindById err:", zap.Error(err))
			return nil, err
		}
	}
	detail := &mysql.QuestionDetail{
		QuestionInfo:               questionInfo,
		QuestionOption:             optionList,
		KnowledgePointQuestionList: questionPointList,
		KnowledgePointList:         knowledgePointList,
	}
	return detail, nil
}

func (q *QuestionService) QuestionUpdateService(c *gin.Context, params *model.QuestionUpdateInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionUpdateService` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//获取问题详情
	questionInfo := &mysql.TQuestion{QuestionId: params.QuestionId}
	questionInfo, err = questionInfo.Find(c, tx)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> The problem does not exist err:", zap.Error(err))
		return err
	}
	questionDetail, err := q.QuestionDetailService(c, &model.QuestionDetailInput{QuestionId: questionInfo.QuestionId})
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> questionInfo.QuestionDetail err:", zap.Error(err))
		return err
	}
	//修改题目信息
	info := questionDetail.QuestionInfo
	info.Name = params.QuestionName
	info.Context = params.Context
	info.Level = params.Level
	info.Answer = params.Answer
	info.Type = params.Type
	info.CreateUser = params.CreateUser
	info.UpdatedAt = time.Now()
	if err = info.Save(c, tx); err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> TQuestion.add err:", zap.Error(err))
		return err
	}
	//修改问题对应知识点编号
	//删除关联
	oldKnowledgeQuestion := &mysql.TKnowledgePointQuestion{QuestionId: params.QuestionId}
	if err = oldKnowledgeQuestion.DeleteById(c, tx); err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> oldKnowledgeQuestion.Delete err:", zap.Error(err))
		return err
	}
	//重新插入
	knpIdList := params.GetKnpIdByModel()
	for _, knp := range knpIdList {
		pointQuestions := &mysql.TKnowledgePointQuestion{QuestionId: params.QuestionId, KnpId: knp}
		if err = pointQuestions.Save(c, tx); err != nil {
			tx.Rollback()
			logger.L.Error("`QuestionUpdateService` -> TKnowledgePointQuestion.Save err:", zap.Error(err))
			return err
		}
	}
	//若是选择题，修改选项表内容
	if params.Type == 1 {
		//删除关联
		oldQuestionOption := &mysql.TQuestionOption{QuestionId: params.QuestionId}
		if err = oldQuestionOption.DeleteById(c, tx); err != nil {
			tx.Rollback()
			logger.L.Error("`QuestionUpdateService` -> oldQuestionOption.Delete err:", zap.Error(err))
			return err
		}
		//重新插入
		for _, option := range params.Option {
			questionOption := &mysql.TQuestionOption{
				QuestionId: questionInfo.QuestionId,
				Context:    option.Context,
				IsAnswer:   option.IsAnswer,
			}
			if err = questionOption.Save(c, tx); err != nil {
				tx.Rollback()
				logger.L.Error("`QuestionAddService` -> TQuestionOption.Save err:", zap.Error(err))
				return err
			}
		}
	}
	tx.Commit()
	return nil
}
