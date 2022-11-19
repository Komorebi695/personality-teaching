package model

type Teacher struct {
	TeacherID   string `gorm:"teacher_id;not null"`
	Name        string `gorm:"name;not null"`
	Password    string `gorm:"password;not null"`
	College     string `gorm:"college;not null"`
	Major       string `gorm:"major;not null"`
	PhoneNumber string `gorm:"phone_number;not null"`
}

func (t Teacher) TableName() string {
	return "t_teacher"
}

type TeacherLoginReq struct {
	UserName string `binding:"required,min=1,max=20" form:"username"`
	Password string `binding:"required,min=1" form:"password"`
}

type SessionValue struct {
	UserID     string `json:"user_id"`
	RoleType   int8   `json:"role_type"`
	CreateTime int64  `json:"create_time"`
}

type TeacherInfoResp struct {
	TeacherID   string `json:"teacher_id"`
	Name        string `json:"name"`
	College     string `json:"college"`
	Major       string `json:"major"`
	PhoneNumber string `json:"phone_number"`
}
