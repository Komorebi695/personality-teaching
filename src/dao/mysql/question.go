package mysql

type QuestionDetail struct {
	QuestionInfo               *TQuestion                 `json:"problem_info" description:"题目信息"`
	KnowledgePointQuestionList []*TKnowledgePointQuestion `json:"knowledge_point_question_list" description:"知识点编号列表"`
	KnowledgePointList         []*TKnowledgePoint         `json:"knowledge_point_list" description:"知识点列表"`
}
