package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
	"time"
)

const SingleChoiceQuestion = 1
const MultipleChoiceQuestions = 2

type QuestionService struct{}

type questionFunc interface {
	QuestionListService(c *gin.Context, params *model.QuestionListInput) (*model.QuestionListOutput, error)
	QuestionDeleteService(c *gin.Context, params *model.QuestionDeleteInput) error
	QuestionAddService(c *gin.Context, params *model.QuestionAddInput) error
	QuestionDetailService(c *gin.Context, params *model.QuestionDetailInput) (*mysql.QuestionDetail, error)
	QuestionUpdateService(c *gin.Context, params *model.QuestionUpdateInput) error
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
		var optionList []*model.QuestionOption
		if listItem.Type == SingleChoiceQuestion || listItem.Type == MultipleChoiceQuestions {
			contextSlice := utils.SplitContext(listItem.QuestionId, listItem.Context)
			// 题干：contextSlice[0]	选项表JSON：contextSlice[1]
			listItem.Context = contextSlice[0]
			if len(contextSlice) == 2 {
				err = json.Unmarshal([]byte(contextSlice[1]), &optionList)
				if err != nil {
					logger.L.Error("`QuestionListService` -> json.Unmarshal err:", zap.Error(err))
					return nil, err
				}
			}
		}
		outItem := model.QuestionListItemOutput{
			QuestionId:   listItem.QuestionId,
			QuestionName: listItem.Name,
			Context:      listItem.Context,
			Option:       optionList,
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
	questionInfo, err = questionInfo.FindOnce(c, tx)
	if err != nil {
		logger.L.Error("`QuestionDeleteService` -> TQuestion.FindOneById err:", zap.Error(err))
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
	if _, err = questionInfo.FindOnce(c, tx); err == nil {
		tx.Rollback()
		logger.L.Error("`QuestionAddService` -> The problem's name already exists:", zap.Error(err))
		return err
	}
	//包装题目信息
	questionId := utils.GenSnowID()
	//若是选择题，选项内容转为JSON插入
	if params.Type == SingleChoiceQuestion || params.Type == MultipleChoiceQuestions {
		//QuestionOptionList	JSON序列化
		optionContext := utils.Obj2Json(params.QuestionOptionList)
		//以生成的questionID后4位取前3位作为分隔符
		splitNum := utils.SplitNum(questionId)
		//context拼接
		params.Context = params.Context + splitNum + optionContext
	}

	questionModel := &mysql.TQuestion{
		//使用雪花ID生成
		QuestionId: questionId,
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
	questionInfo, err = questionInfo.FindOnce(c, tx)
	if err != nil {
		logger.L.Error("`QuestionDetailService` -> questionInfo.FindOneById err:", zap.Error(err))
		return nil, err
	}
	//题目对应的知识点编号表
	questionPointSearch := &mysql.TKnowledgePointQuestion{QuestionId: questionInfo.QuestionId}
	questionPointList, err := questionPointSearch.Find(c, tx, questionPointSearch)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.L.Error("`QuestionDetailService` -> TKnowledgePointQuestion.FindOneById err:", zap.Error(err))
		return nil, err
	}
	//根据编号表查询知识点列表
	var knowledgePointList []*mysql.TKnowledgePoint
	for _, point := range questionPointList {
		pointSearch := &mysql.TKnowledgePoint{KnpId: point.KnpId}
		pointSearch, err = pointSearch.FindOneById(c, tx)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.L.Error("`QuestionDetailService` -> TKnowledgePoint.FindOneById err:", zap.Error(err))
			return nil, err
		}
		knowledgePointList = append(knowledgePointList, pointSearch)
	}
	//若类型是选择题，获取选项结构体
	var optionList []*model.QuestionOption
	if questionInfo.Type == SingleChoiceQuestion || questionInfo.Type == MultipleChoiceQuestions {
		contextSlice := utils.SplitContext(questionInfo.QuestionId, questionInfo.Context)
		// 题干：contextSlice[0]	选项表JSON：contextSlice[1]
		questionInfo.Context = contextSlice[0]
		if len(contextSlice) == 2 {
			err = json.Unmarshal([]byte(contextSlice[1]), &optionList)
			if err != nil {
				logger.L.Error("`QuestionDetailService` -> json.Unmarshal err:", zap.Error(err))
				return nil, err
			}
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
	questionInfo, err = questionInfo.FindOnce(c, tx)
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
	info := questionDetail.QuestionInfo
	//判断题目类型是否为选择题
	if params.Type == SingleChoiceQuestion || params.Type == MultipleChoiceQuestions {
		//QuestionOptionList	JSON序列化
		optionContext := utils.Obj2Json(params.Option)
		//以生成的questionID后3位作为分隔符
		splitNum := utils.SplitNum(info.QuestionId)
		//context拼接
		params.Context = params.Context + splitNum + optionContext
	}
	//修改题目信息
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

	tx.Commit()
	return nil
}
