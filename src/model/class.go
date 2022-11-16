package model

type Class struct {
	ClassID string `gorm:"column:class_id" json:"class_id"`
	BaseClassInfo
}

func (c Class) TableName() string {
	return "t_class"
}

type BaseClassInfo struct {
	Name    string `json:"name" form:"name" binding:"required" gorm:"column:name"`
	College string `json:"college" form:"college" binding:"required" gorm:"column:college"`
	Major   string `json:"major" form:"major" binding:"required" gorm:"column:major"`
}

type ClassAddReq struct {
	BaseClassInfo
}

type ClassUpdateReq struct {
	ClassID string `form:"class_id" binding:"required"`
	BaseClassInfo
}

type ClassDeleteReq struct {
	ClassID string `form:"class_id" binding:"required"`
}

type ClassInfoReq struct {
	ClassID string `form:"class_id" binding:"required"`
}

type ClassListReq struct {
	PageNum  int `form:"page_num" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}
