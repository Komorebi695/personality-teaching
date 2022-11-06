package logic

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

type KnowledgePointService struct{}

type knowledgePointFunc interface {
	KnowledgePointListService(c *gin.Context, params *model.KnowledgePointListInput) (*model.KnowledgePointListOutput, error)
	KnowledgePointDeleteService(c *gin.Context, params *model.KnowledgePointDeleteInput) error
	KnowledgePointAddService(c *gin.Context, params *model.KnowledgePointAddInput) error
	KnowledgePointDetailService(c *gin.Context, params *model.KnowledgePointDetailInput) (*mysql.TKnowledgePoint, error)
}

var _ knowledgePointFunc = &KnowledgePointService{}

func NewKnowledgePointService() *KnowledgePointService {
	return &KnowledgePointService{}
}

// KnowledgePointListService 知识点列表查询
func (q *KnowledgePointService) KnowledgePointListService(c *gin.Context, params *model.KnowledgePointListInput) (*model.KnowledgePointListOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointListService` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//从db中分页读取基本信息
	knowledgePointInfo := &mysql.TKnowledgePoint{}
	list, total, err := knowledgePointInfo.PageList(c, tx, params)
	if err != nil {
		logger.L.Error("`KnowledgePointListService` -> knowledgePointInfo.PageList err:", zap.Error(err))
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

// KnowledgePointDeleteService 知识点删除
func (q *KnowledgePointService) KnowledgePointDeleteService(c *gin.Context, params *model.KnowledgePointDeleteInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointDeleteService` -> get pool err:", zap.Error(err))
		return err
	}
	//读取基本信息
	knowledgePointInfo := &mysql.TKnowledgePoint{KnpId: params.KnpId}
	knowledgePointInfo, err = knowledgePointInfo.FindOneById(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDeleteService` -> TKnowledgePoint.FindOneById err:", zap.Error(err))
		return err
	}
	err = knowledgePointInfo.Delete(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDeleteService` -> TKnowledgePoint.Delete err:", zap.Error(err))
		return err
	}
	return nil
}

// KnowledgePointAddService 知识点添加
func (q *KnowledgePointService) KnowledgePointAddService(c *gin.Context, params *model.KnowledgePointAddInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointAddService` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//判断知识点是否重复插入
	knowledgePointInfo := &mysql.TKnowledgePoint{Name: params.Name}
	if _, err = knowledgePointInfo.FindByName(c, tx); err == nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointAddService` -> The KnowledgePoint already exists:", zap.Error(err))
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
		logger.L.Error("`KnowledgePointAddService` -> knowledgePointModel.Save err:", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

// KnowledgePointDetailService 知识点详情
func (q *KnowledgePointService) KnowledgePointDetailService(c *gin.Context, params *model.KnowledgePointDetailInput) (*mysql.TKnowledgePoint, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointDetailService` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//获取知识点详情
	knowledgePointInfo := &mysql.TKnowledgePoint{KnpId: params.KnpId}
	knowledgePointInfo, err = knowledgePointInfo.FindOneById(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointDetailService` -> knowledgePointInfo.FindOneById err:", zap.Error(err))
		return nil, err
	}
	return knowledgePointInfo, nil
}

// KnowledgePointUpdateService 知识点修改
func (q *KnowledgePointService) KnowledgePointUpdateService(c *gin.Context, params *model.KnowledgePointUpdateInput) error {
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
