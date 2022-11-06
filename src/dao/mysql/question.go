package mysql

import "personality-teaching/src/model"

type QuestionDetail struct {
	QuestionInfo               *TQuestion                 `json:"problem_info" description:"题目信息"`
	QuestionOption             []*model.QuestionOption    `json:"question_option" comment:"选项信息"`
	KnowledgePointQuestionList []*TKnowledgePointQuestion `json:"knowledge_point_question_list" description:"知识点编号列表"`
	KnowledgePointList         []*TKnowledgePoint         `json:"knowledge_point_list" description:"知识点列表"`
}
