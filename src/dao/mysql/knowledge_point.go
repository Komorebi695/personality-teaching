package mysql

import "personality-teaching/src/model"

type KnowledgePointDetail struct {
	Info                    *TKnowledgePoint                 `json:"info" description:"知识点信息"`
	Children                []*TKnowledgePoint               `json:"children" description:"子知识点列表信息"`
	KnowledgeConnectionList []*model.KnowledgeConnectionItem `json:"knowledge_connection_list" description:"知识点联系列表"`
}
