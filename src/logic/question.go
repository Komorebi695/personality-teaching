package logic

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
	"strings"
	"time"
)

const SingleChoiceQuestion = 1
const MultipleChoiceQuestions = 2

type QuestionService struct {
	knpQuestionArticle *mysql.KnowledgePointQuestionMySQL
	questionArticle    *mysql.QuestionMySQL
}

type questionFunc interface {
	QuestionListService(c *gin.Context, params *model.QuestionListInput) (*model.QuestionListOutput, error)
	QuestionDeleteService(c *gin.Context, params *model.QuestionDeleteInput) error
	QuestionAddService(c *gin.Context, params *model.QuestionAddInput) error
	QuestionDetailService(c *gin.Context, params *model.QuestionDetailInput) (*model.QuestionDetail, error)
	QuestionUpdateService(c *gin.Context, params *model.QuestionUpdateInput) error
}

var _ questionFunc = &QuestionService{}

func NewQuestionService() *QuestionService {
	return &QuestionService{
		knpQuestionArticle: mysql.NewKnowledgePointQuestionMySQL(),
		questionArticle:    mysql.NewQuestionMySQL(),
	}
}

func (q *QuestionService) QuestionListService(c *gin.Context, params *model.QuestionListInput) (*model.QuestionListOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionListService` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//从db中分页读取基本信息
	list, total, err := q.questionArticle.PageList(c, tx, params)
	if err != nil {
		logger.L.Error("`QuestionListService` -> questionInfo.PageList err:", zap.Error(err))
		return nil, err
	}
	//格式化输出信息
	var outList []model.QuestionListItemOutput
	for _, listItem := range list {

		//拆分问题选项和问题内容
		context, optionList, err := OptionSpit(listItem)
		if err != nil {
			return nil, err
		}
		//拆分问题答案和答案解析
		answer, answerContext, err := AnswerSpit(listItem)
		if err != nil {
			return nil, err
		}
		outItem := model.QuestionListItemOutput{
			QuestionBase: model.QuestionBase{
				Name:    listItem.Name,
				Level:   listItem.Level,
				Type:    listItem.Type,
				Context: context,
				Answer:  answer,
			},
			QuestionId:    listItem.QuestionId,
			Option:        optionList,
			AnswerContext: answerContext,
			CreateUser:    listItem.CreateUser,
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
	questionInfo, err := q.questionArticle.FindOnce(c, tx, params.QuestionId)
	if err != nil {
		logger.L.Error("`QuestionDeleteService` -> Question.FindOneById err:", zap.Error(err))
		return err
	}
	if err = q.questionArticle.Delete(c, tx, questionInfo.ID); err != nil {
		logger.L.Error("`QuestionDeleteService` -> Question.Save err:", zap.Error(err))
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
	if _, err = q.questionArticle.FindOnce(c, tx, params.QuestionBase.Name); err == nil {
		tx.Rollback()
		logger.L.Error("`QuestionAddService` -> The problem's name already exists:", zap.Error(err))
		return errors.New("the problem's name already exists")
	}
	//包装题目信息
	//使用雪花ID生成questionId
	questionId := utils.GenSnowID()
	//以生成的questionID后4位取前3位作为分隔符
	splitNum, err := utils.SplitNum(questionId)
	if err != nil {
		return err
	}
	//若是选择题，选项内容转为JSON插入
	if params.Type == SingleChoiceQuestion || params.Type == MultipleChoiceQuestions {
		context, err := OptionSplice(params.QuestionOptionList, params.Context, splitNum)
		if err != nil {
			return err
		}
		params.Context = context
	}
	//包装答案与答案解析
	answer, err := AnswerSplice(params.Answer, splitNum, params.AnswerContext)
	if err != nil {
		return err
	}
	params.Answer = answer
	questionModel := &model.Question{
		QuestionId: questionId,
		QuestionBase: model.QuestionBase{
			Name:    params.Name,
			Level:   params.Level,
			Type:    params.Type,
			Context: params.Context,
			Answer:  answer,
		},
		CreateUser: params.CreateUser,
	}
	if err = q.questionArticle.Save(c, tx, questionModel); err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionAddService` -> questionModel.Save err:", zap.Error(err))
		return err
	}

	//循环依次插入知识点题目关联表
	knpIdList := params.GetKnpIdByModel()
	for _, s := range knpIdList {
		knowledgePointQuestion := &model.KnowledgePointQuestion{
			KnpId:      s,
			QuestionId: questionModel.QuestionId,
		}
		if err = q.knpQuestionArticle.Save(c, tx, knowledgePointQuestion); err != nil {
			tx.Rollback()
			logger.L.Error("`QuestionAddService` -> knowledgePointQuestion.Save err:", zap.Error(err))
			return err
		}
	}
	tx.Commit()
	return nil
}

func (q *QuestionService) QuestionDetailService(c *gin.Context, params *model.QuestionDetailInput) (*model.QuestionDetail, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`QuestionDetailService` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//获取问题详情
	question, err := q.questionArticle.FindOnce(c, tx, params.QuestionId)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.L.Error("`QuestionDetailService` -> questionInfo.FindOneById err:", zap.Error(err))
		}
		return nil, err
	}
	//题目对应的知识点编号表
	questionPointList, err := q.knpQuestionArticle.Find(c, tx, question.QuestionId)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.L.Error("`QuestionDetailService` -> KnowledgePointQuestion.FindOneById err:", zap.Error(err))
		return nil, err
	}
	//根据编号表查询知识点列表
	var knowledgePointList []*model.KnowledgePoint
	for _, point := range questionPointList {
		pointSearch, err := mysql.NewKnowledgePointMySQL().FindOneById(c, tx, point.KnpId)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.L.Error("`QuestionDetailService` -> TKnowledgePoint.FindOneById err:", zap.Error(err))
			return nil, err
		}
		knowledgePointList = append(knowledgePointList, pointSearch)
	}
	//拆分问题选项和问题内容
	context, optionList, err := OptionSpit(*question)
	if err != nil {
		return nil, err
	}
	question.Context = context
	//拆分问题答案和答案解析
	answer, answerContext, err := AnswerSpit(*question)
	if err != nil {
		return nil, err
	}
	question.Answer = answer
	//拼接题目信息
	questionInfo := &model.QuestionInfo{
		Question:       *question,
		AnswerContext:  answerContext,
		QuestionOption: optionList,
	}
	//拼接题目详情
	detail := &model.QuestionDetail{
		QuestionInfo:       *questionInfo,
		KnowledgePointList: knowledgePointList,
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
	question, err := q.questionArticle.FindOnce(c, tx, params.QuestionId)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> The problem does not exist err:", zap.Error(err))
		return err
	}
	questionDetail, err := q.QuestionDetailService(c, &model.QuestionDetailInput{QuestionId: question.QuestionId})
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> questionInfo.QuestionDetail err:", zap.Error(err))
		return err
	}
	info := questionDetail.QuestionInfo
	//以生成的questionID后3位作为分隔符
	splitNum, err := utils.SplitNum(info.QuestionId)
	if err != nil {
		return err
	}
	//判断题目类型是否为选择题
	//若是选择题，选项内容转为JSON插入context
	if params.Type == SingleChoiceQuestion || params.Type == MultipleChoiceQuestions {
		context, err := OptionSplice(params.Option, params.Context, splitNum)
		if err != nil {
			return err
		}
		params.Context = context
	}

	//包装答案与答案解析
	answer, err := AnswerSplice(params.Answer, splitNum, params.AnswerContext)
	if err != nil {
		return err
	}
	params.Answer = answer
	//修改题目信息
	info.Name = params.Name
	info.Context = params.Context
	info.Level = params.Level
	info.Answer = params.Answer
	info.Type = params.Type
	info.CreateUser = params.CreateUser
	info.UpdatedAt = time.Now()
	if err = q.questionArticle.Save(c, tx, &info.Question); err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> Question.add err:", zap.Error(err))
		return err
	}
	//修改问题对应知识点编号
	//删除关联
	if err = q.knpQuestionArticle.DeleteAllById(c, tx, params.QuestionId); err != nil {
		tx.Rollback()
		logger.L.Error("`QuestionUpdateService` -> oldKnowledgeQuestion.Delete err:", zap.Error(err))
		return err
	}
	//重新插入
	knpIdList := params.GetKnpIdByModel()
	for _, knp := range knpIdList {
		pointQuestions := &model.KnowledgePointQuestion{QuestionId: params.QuestionId, KnpId: knp}
		if err = q.knpQuestionArticle.Save(c, tx, pointQuestions); err != nil {
			tx.Rollback()
			logger.L.Error("`QuestionUpdateService` -> KnowledgePointQuestion.Save err:", zap.Error(err))
			return err
		}
	}
	tx.Commit()
	return nil
}

// OptionSpit 拆分题目选项和题目内容
func OptionSpit(question model.Question) (string, []*model.QuestionOption, error) {
	var optionList []*model.QuestionOption
	if question.Type == SingleChoiceQuestion || question.Type == MultipleChoiceQuestions {
		contextSlice, err := utils.SplitContext(question.QuestionId, question.Context)
		if err != nil {
			return "", nil, err
		}
		// 题干：contextSlice[0]	选项表JSON：contextSlice[1]
		context := contextSlice[0]
		if len(contextSlice) == 2 {
			err := json.Unmarshal([]byte(contextSlice[1]), &optionList)
			if err != nil {
				return "", nil, errors.New("json.Unmarshal err")
			}
		}
		return context, optionList, nil
	}
	return question.Context, nil, nil
}

// OptionSplice 拼接选项与内容
func OptionSplice(optionList []*model.QuestionOption, context string, splitNum string) (string, error) {
	//选项内容序列化为JSON插入
	optionContext, err := utils.Obj2Json(optionList)
	if err != nil {
		return "", err
	}
	//context拼接
	var c strings.Builder
	c.WriteString(context)
	c.WriteString(splitNum)
	c.WriteString(optionContext)
	return c.String(), nil
}

// AnswerSpit 拆分问题答案和答案解析
func AnswerSpit(question model.Question) (string, string, error) {
	answerContextSlice, err := utils.SplitContext(question.QuestionId, question.Answer)
	if err != nil {
		return "", "", err
	}
	// 答案：contextSlice[0]	解析：contextSlice[1]
	var answer, answerContext string
	if len(answerContextSlice) == 2 {
		answer = answerContextSlice[0]
		answerContext = answerContextSlice[1]
	}
	return answer, answerContext, nil
}

// AnswerSplice 拼接答案与答案解析
func AnswerSplice(answer, splitNum, answerContext string) (string, error) {
	//answer拼接
	var a strings.Builder
	a.WriteString(answer)
	a.WriteString(splitNum)
	a.WriteString(answerContext)
	return a.String(), nil
}
