package logic

import (
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

type studentFunc interface {
	CreateStudent(req model.CreateStudentReq) (model.CreateStudentResp, error)
	UpdateClassID(req model.AddStudentToClassReq) (model.AddStudentClassResp, error)
	GetStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, error)
	RemoveStudentClass(studentID string) error
	CheckStudentClass(studentID string, classID string) (bool, error)
}

var _ studentFunc = &StudentService{}

type StudentService struct{}

func NewStudentService() *StudentService {
	return &StudentService{}
}

func (s *StudentService) CreateStudent(req model.CreateStudentReq) (model.CreateStudentResp, error) {
	student := model.Student{
		BaseStudentInfo: req.BaseStudentInfo,
		StudentID:       utils.GenSnowID(),
		Password:        utils.GetDefaultPassWord(),
		ClassID:         utils.EmptyClassID,
	}
	if err := mysql.NewStudentMySQL().InsertStudent(student); err != nil {
		return model.CreateStudentResp{}, err
	}
	return model.CreateStudentResp{
		StudentID:       student.StudentID,
		BaseStudentInfo: student.BaseStudentInfo,
	}, nil
}

func (s *StudentService) UpdateClassID(req model.AddStudentToClassReq) (model.AddStudentClassResp, error) {
	if err := mysql.NewStudentMySQL().UpdateClassID(req.StudentID, req.ClassID); err != nil {
		return model.AddStudentClassResp{}, err
	}
	student, err := mysql.NewStudentMySQL().QueryStudent(req.StudentID)
	if err != nil {
		return model.AddStudentClassResp{}, err
	}
	return model.AddStudentClassResp{
		StudentID:       student.StudentID,
		ClassID:         student.ClassID,
		BaseStudentInfo: student.BaseStudentInfo,
	}, nil
}

func (s *StudentService) GetStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, error) {
	students, err := mysql.NewStudentMySQL().QueryStudentsInClass(req)
	if err != nil {
		return []model.ClassStudentListResp{}, err
	}
	return students, nil
}

func (s *StudentService) RemoveStudentClass(studentID string) error {
	return mysql.NewStudentMySQL().UpdateClassID(studentID, utils.EmptyClassID)
}

func (s *StudentService) CheckStudentClass(studentID string, classID string) (bool, error) {
	return mysql.NewStudentMySQL().CheckStudentClass(studentID, classID)
}
