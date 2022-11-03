package model

type TypeQueryInput struct {
	Id int8 `json:"id" form:"id" comment:"题目类型Id" binding:"required"`
}

// TypeQueryOutput 题目类型查找输出
type TypeQueryOutput struct {
	Id   int8   `json:"id" form:"id" comment:"题目类型Id，1：选择题，2：填空题，3：问答题"`
	Name string `json:"type_name" form:"type_name" comment:"题目类型名称、选择题、填空题、问答题"`
}
