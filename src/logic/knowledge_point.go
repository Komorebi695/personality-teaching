package logic

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/logger"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

type KnowledgePointService struct {
	knpArticle     *mysql.KnowledgePointMySQL
	knpConnArticle *mysql.KnowledgeConnectionMySQL
}

type knowledgePointFunc interface {
	KnowledgePointList(c *gin.Context, params *model.KnowledgePointListInput) (*model.KnowledgePointListOutput, error)
	KnowledgePointOneStageList(c *gin.Context) (*model.KnpOneStageListOutput, error)
	KnowledgePointDelete(c *gin.Context, params *model.KnowledgePointDeleteInput) error
	KnowledgePointAdd(c *gin.Context, params *model.KnowledgePointAddInput) error
	KnowledgePointDetail(c *gin.Context, params *model.KnowledgePointDetailInput) (*model.KnowledgePointDetail, error)
	KnowledgePointUpdate(c *gin.Context, params *model.KnowledgePointUpdateInput) error
	PointConnectionUpdate(c *gin.Context, params *model.KnpConnectionUpdateInput) error
}

var _ knowledgePointFunc = &KnowledgePointService{}

func NewKnowledgePointService() *KnowledgePointService {
	return &KnowledgePointService{
		mysql.NewKnowledgePointMySQL(),
		mysql.NewKnowledgeConnectionMySQL(),
	}
}

// KnowledgePointList 知识点列表查询
func (t *KnowledgePointService) KnowledgePointList(c *gin.Context, params *model.KnowledgePointListInput) (*model.KnowledgePointListOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//从db中分页读取基本信息
	list, total, err := t.knpArticle.PageList(c, tx, params)
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> knowledgePointInfo.PageList err:", zap.Error(err))
		return nil, err
	}
	//格式化输出信息
	out := &model.KnowledgePointListOutput{
		Total: total,
		List:  list,
	}
	return out, nil
}

// KnowledgePointOneStageList 一级知识点列表查询
func (t *KnowledgePointService) KnowledgePointOneStageList(c *gin.Context) (*model.KnpOneStageListOutput, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//从db中分页读取基本信息
	list, err := t.knpArticle.PageListOneStage(c, tx)
	if err != nil {
		logger.L.Error("`KnowledgePointList` -> knowledgePointInfo.PageList err:", zap.Error(err))
		return nil, err
	}
	//格式化输出信息
	out := &model.KnpOneStageListOutput{
		List: list,
	}
	return out, nil
}

// KnowledgePointDelete 知识点删除（根据knpId删除当前知识点以及所有子节点）
func (t *KnowledgePointService) KnowledgePointDelete(c *gin.Context, params *model.KnowledgePointDeleteInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointDelete` -> get pool err:", zap.Error(err))
		return err
	}
	//开启事务
	tx = tx.Begin()
	err = deleteKnpAndChild(c, tx, params.KnpId, t)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointDelete` -> TKnowledgePoint.deleteKnpAndChild err:", zap.Error(err))
		return err
	}
	//提交事务
	tx.Commit()
	return nil
}

//根据knpId删除当前知识点以及所有子节点
func deleteKnpAndChild(c *gin.Context, tx *gorm.DB, knpId string, t *KnowledgePointService) error {
	//查询知识点信息
	knowledgePointInfo, err := t.knpArticle.FindOneById(c, tx, knpId)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`deleteKnpAndChild` -> TKnowledgePoint.FindOneById err:", zap.Error(err))
		return err
	}
	//查询该知识点是否存在子知识点
	children, err := t.knpArticle.FindKnowledgeChildren(c, tx, knpId)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`deleteKnpAndChild` -> knowledgePointInfo.FindKnowledgeChildren err:", zap.Error(err))
		return err
	}
	// 若存在子知识点，遍历删除
	if len(children) != 0 {
		for _, child := range children {
			err := deleteKnpAndChild(c, tx, child.KnpId, t)
			if err != nil {
				tx.Rollback()
				logger.L.Error("`deleteKnpAndChild` -> deleteKnpAndChild.Delete Child err:", zap.Error(err))
				return err
			}
		}
	}
	err = t.knpArticle.Delete(c, tx, knowledgePointInfo.Id)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`deleteKnpAndChild` -> deleteKnpAndChild.Delete Knp err:", zap.Error(err))
		return err
	}
	return nil
}

// KnowledgePointAdd 知识点添加
func (t *KnowledgePointService) KnowledgePointAdd(c *gin.Context, params *model.KnowledgePointAddInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointAdd` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//判断知识点是否重复插入
	if _, err = t.knpArticle.FindByName(c, tx, params.Name); err == nil {
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
	knowledgePointModel := &model.KnowledgePoint{
		KnowledgePointBase: model.KnowledgePointBase{
			KnpId:       knpId,
			ParentKnpId: params.ParentKnpId,
		},
		KnowledgePointInfo: model.KnowledgePointInfo{
			Name:    params.Name,
			Level:   params.Level,
			Context: params.Context,
		},
		CreateUser: params.CreateUser,
	}
	if err = t.knpArticle.Save(c, tx, knowledgePointModel); err != nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointAdd` -> knowledgePointModel.Save err:", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

// KnowledgePointDetail 知识点详情
func (t *KnowledgePointService) KnowledgePointDetail(c *gin.Context, params *model.KnowledgePointDetailInput) (*model.KnowledgePointDetail, error) {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointDetail` -> get pool err:", zap.Error(err))
		return nil, err
	}
	//获取知识点详情
	//知识点信息
	knowledgePointInfo, err := t.knpArticle.FindOneById(c, tx, params.KnpId)
	if err != nil {
		logger.L.Error("`KnowledgePointDetail` -> knowledgePointInfo.FindOneById err:", zap.Error(err))
		return nil, err
	}
	//子知识点列表
	children, err := t.knpArticle.FindKnowledgeChildren(c, tx, params.KnpId)
	if err != nil {
		logger.L.Error("`KnowledgePointDetail` -> knowledgePointInfo.FindKnowledgeChildren err:", zap.Error(err))
		return nil, err
	}

	//知识点联系列表
	connectionList, err := t.knpConnArticle.QueryNameById(c, tx, params.KnpId)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.L.Error("`KnowledgePointDetail` -> connectionInfo.Find err:", zap.Error(err))
		return nil, err
	}
	out := &model.KnowledgePointDetail{
		Info:                    knowledgePointInfo,
		Children:                children,
		KnowledgeConnectionList: connectionList,
	}
	return out, nil
}

// KnowledgePointUpdate 知识点修改
func (t *KnowledgePointService) KnowledgePointUpdate(c *gin.Context, params *model.KnowledgePointUpdateInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`KnowledgePointUpdateService` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//获取知识点详情
	knowledgePointInfo, err := t.knpArticle.FindOneById(c, tx, params.KnpId)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointUpdateService` -> The knowledgePoint does not exist err:", zap.Error(err))
		return err
	}
	//修改知识点信息
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
	if err = t.knpArticle.Save(c, tx, info); err != nil {
		tx.Rollback()
		logger.L.Error("`KnowledgePointUpdateService` -> TKnowledgePoint.save err:", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

// PointConnectionUpdate 修改知识点联系
func (t *KnowledgePointService) PointConnectionUpdate(c *gin.Context, params *model.KnpConnectionUpdateInput) error {
	tx, err := mysql.GetGormPool()
	if err != nil {
		logger.L.Error("`PointConnectionUpdate Service` -> get pool err:", zap.Error(err))
		return err
	}
	tx = tx.Begin()
	//获取知识点联系列表
	//修改知识点联系
	//全删了，重新插入
	err = t.knpConnArticle.DeleteById(c, tx, params.KnpId)
	if err != nil {
		tx.Rollback()
		logger.L.Error("`PointConnectionUpdate Service` -> connectionInfo.DeleteAllByKnpId:", zap.Error(err))
		return err
	}
	pKnpIdList := params.GetKnpIdByModel()
	for _, pKnpId := range pKnpIdList {
		Item := &model.KnowledgeConnection{
			KnpId:  params.KnpId,
			PKnpId: pKnpId,
		}
		if err = t.knpConnArticle.Save(c, tx, Item); err != nil {
			tx.Rollback()
			logger.L.Error("`PointConnectionUpdate Service` -> PointConnectionUpdate.save err:", zap.Error(err))
			return err
		}
	}

	tx.Commit()
	return nil
}
