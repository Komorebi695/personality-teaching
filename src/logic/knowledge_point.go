package logic

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

type KnowledgePointService struct{}

type knowledgePointFunc interface {
	KnowledgePointList(c *gin.Context, params *model.KnowledgePointListInput) (*model.KnowledgePointListOutput, error)
	KnowledgePointOneStageList(c *gin.Context) (*model.KnowledgePointOneStageListOutput, error)
	KnowledgePointDelete(c *gin.Context, params *model.KnowledgePointDeleteInput) error
	KnowledgePointAdd(c *gin.Context, params *model.KnowledgePointAddInput) error
	KnowledgePointDetail(c *gin.Context, params *model.KnowledgePointDetailInput) (*mysql.KnowledgePointDetail, error)
}

var _ knowledgePointFunc = &KnowledgePointService{}

func NewKnowledgePointService() *KnowledgePointService {
	return &KnowledgePointService{}
}

// KnowledgePointList 知识点列表查询
func (q *KnowledgePointService) KnowledgePointList(c *gin.Context, params *model.KnowledgePointListInput) (*model.KnowledgePointListOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//从db中分页读取基本信息
	knowledgePointInfo := &mysql.TKnowledgePoint{}
	list, total, err := knowledgePointInfo.PageList(c, tx, params)
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> knowledgePointInfo.PageList err:", zap.Error(err))
		return nil, err
	}
	//格式化输出信息
	var outList []model.KnowledgePointListItemOutput
	for _, listItem := range list {
		outItem := model.KnowledgePointListItemOutput{
			KnpId:       listItem.KnpId,
			ParentKnpId: listItem.ParentKnpId,
			Name:        listItem.Name,
			Level:       listItem.Level,
			Context:     listItem.Context,
			CreateUser:  listItem.CreateUser,
		}
		outList = append(outList, outItem)
	}
	out := &model.KnowledgePointListOutput{
		Total: total,
		List:  outList,
	}
	return out, nil
}

// KnowledgePointOneStageList 一级知识点列表查询
func (q *KnowledgePointService) KnowledgePointOneStageList(c *gin.Context) (*model.KnowledgePointOneStageListOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//从db中分页读取基本信息
	knowledgePointInfo := &mysql.TKnowledgePoint{}
	list, err := knowledgePointInfo.PageListOneStage(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> knowledgePointInfo.PageList err:", zap.Error(err))
		return nil, err
	}
	//格式化输出信息
	var outList []model.KnpOneStageListItemOutput
	for _, listItem := range list {
		outItem := model.KnpOneStageListItemOutput{
			KnpId:       listItem.KnpId,
			ParentKnpId: listItem.ParentKnpId,
			Name:        listItem.Name,
			Level:       listItem.Level,
			Context:     listItem.Context,
		}
		outList = append(outList, outItem)
	}
	out := &model.KnowledgePointOneStageListOutput{
		List: outList,
	}
	return out, nil
}

// KnowledgePointDelete 知识点删除
func (q *KnowledgePointService) KnowledgePointDelete(c *gin.Context, params *model.KnowledgePointDeleteInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointDelete` -> get pool err:", zap.Error(err))
		return err
	}
	//查询基本信息
	knowledgePointInfo := &mysql.TKnowledgePoint{KnpId: params.KnpId}
	knowledgePointInfo, err = knowledgePointInfo.FindOneById(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDelete` -> TKnowledgePoint.FindOneById err:", zap.Error(err))
		return err
	}
	//查询该知识点是否存在子知识点
	children, err := knowledgePointInfo.FindKnowledgeChildren(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDelete` -> knowledgePointInfo.FindKnowledgeChildren err:", zap.Error(err))
		return err
	}
	// 若存在子知识点，删除失败返回
	if len(children) != 0 {
		err = errors.New("child node exists err")
		logger.L.Error("`KnowledgePointDelete` -> Child KnowledgePoint exists err:", zap.Error(err))
		return err
	}
	err = knowledgePointInfo.Delete(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDelete` -> TKnowledgePoint.Delete err:", zap.Error(err))
		return err
	}
	return nil
}

// KnowledgePointAdd 知识点添加
func (q *KnowledgePointService) KnowledgePointAdd(c *gin.Context, params *model.KnowledgePointAddInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointAdd` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//判断知识点是否重复插入
	knowledgePointInfo := &mysql.TKnowledgePoint{Name: params.Name}
	if _, err = knowledgePointInfo.FindByName(c, tx); err == nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointAdd` -> The KnowledgePoint already exists:", zap.Error(err))
		return err
	}
	knpId := utils.GenSnowID()
	// 若父知识点为空，则默认指向自己
	if params.ParentKnpId == "" {
		params.ParentKnpId = knpId
	}
	//包装知识点信息
	knowledgePointModel := &mysql.TKnowledgePoint{
		KnpId:       knpId,
		Name:        params.Name,
		ParentKnpId: params.ParentKnpId,
		Level:       params.Level,
		Context:     params.Context,
		CreateUser:  params.CreateUser,
	}
	if err = knowledgePointModel.Save(c, tx); err != nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointAdd` -> knowledgePointModel.Save err:", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

// KnowledgePointDetail 知识点详情
func (q *KnowledgePointService) KnowledgePointDetail(c *gin.Context, params *model.KnowledgePointDetailInput) (*mysql.KnowledgePointDetail, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointDetail` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//获取知识点详情
	knowledgePointInfo := &mysql.TKnowledgePoint{KnpId: params.KnpId}
	knowledgePointInfo, err = knowledgePointInfo.FindOneById(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDetail` -> knowledgePointInfo.FindOneById err:", zap.Error(err))
		return nil, err
	}
	children, err := knowledgePointInfo.FindKnowledgeChildren(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDetail` -> knowledgePointInfo.FindKnowledgeChildren err:", zap.Error(err))
		return nil, err
	}
	out := &mysql.KnowledgePointDetail{
		Info:     knowledgePointInfo,
		Children: children,
	}
	return out, nil
}

// KnowledgePointUpdate 知识点修改
func (q *KnowledgePointService) KnowledgePointUpdate(c *gin.Context, params *model.KnowledgePointUpdateInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointUpdateService` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//获取知识点详情
	knowledgePointInfo := &mysql.TKnowledgePoint{KnpId: params.KnpId}
	knowledgePointInfo, err = knowledgePointInfo.FindOneById(c, tx)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointUpdateService` -> The knowledgePoint does not exist err:", zap.Error(err))
		return err
	}
	//修改题目信息

	// 若父知识点为空，则默认指向自己
	if params.ParentKnpId == "" {
		params.ParentKnpId = params.KnpId
	}
	info := knowledgePointInfo
	info.Name = params.Name
	info.Context = params.Context
	info.Level = params.Level
	info.ParentKnpId = params.ParentKnpId
	info.CreateUser = params.CreateUser
	if err = info.Save(c, tx); err != nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointUpdateService` -> TKnowledgePoint.save err:", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}
