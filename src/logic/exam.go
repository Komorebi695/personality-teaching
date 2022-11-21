package logic

import (
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

const (
	one int8 = 1 // 常量 1
)

type ExamService struct{}

func NewExamService() *ExamService {
	return &ExamService{}
}

// SearchExam ,新增试卷逻辑处理部分
//Param：
// text：搜索文本
// teacherID：老师编号
//Return value：
// 试卷的信息和错误
func (ec *ExamService) SearchExam(text string, teacherID string) (model.ExamListResp, error) {
	return mysql.NewExamMysql().Query(text, teacherID)
}

// Add ,新增试卷逻辑处理部分
//Param：
// teacherID：老师编号
// req：试卷的基本信息()
//Return value：
// 试卷的信息和错误
func (ec *ExamService) Add(teacherID string, req model.ExamAddReq) (model.Exam, error) {
	exam := model.Exam{
		ExamID:          utils.GenSnowID(),
		CreateTeacherID: teacherID,
		UpdateTime:      utils.CurrentTime(),
		CreateTime:      utils.CurrentTime(),
		BaseExamInfo:    req.BaseExamInfo,
	}
	if err := mysql.NewExamMysql().Insert(exam); err != nil {
		return model.Exam{}, err
	}
	return exam, nil
}

// Update ,更新试卷
//Param：
// req：试卷的基本信息()
//Return value：
// 错误信息
func (ec *ExamService) Update(req model.ExamUpdateReq) error {
	exam := model.Exam{
		ExamID:       req.ExamID,
		BaseExamInfo: req.BaseExamInfo,
		UpdateTime:   utils.CurrentTime(),
	}
	return mysql.NewExamMysql().UpdateExam(exam)
}

// Delete ,删除试卷
//Param：
// exam_id：试卷编号
//Return value：
// 错误信息
func (ec *ExamService) Delete(req model.ExamDeleteReq) error {
	return mysql.NewExamMysql().DeleteExam(req.ExamID)
}

// List ,获取当前老师的所有试卷
//Param：
// teacherID：老师编号
// req: 分页参数
//Return value：
// 所有试卷信息和错误信息
func (ec *ExamService) List(teacherID string, req model.PagingReq) (model.ExamListResp, error) {
	offset := (req.Page - int(one)) * req.PageSize
	return mysql.NewExamMysql().QueryExamList(teacherID, offset, req.PageSize)
}

// Details ,获取试卷详细信息
//Param：
// examID：试卷编号
//Return value：
// 试卷信息和错误信息
func (ec *ExamService) Details(examID string) (model.ExamDetailResp, error) {
	return mysql.NewExamMysql().QueryExam(examID)
}

// SendPerson ,按个人发放卷
//Param：
// StudentID：学生编号
//Return value：
// 错误信息
func (ec *ExamService) SendPerson(req model.SendPersonReq) error {
	se := model.StudentExam{
		SendReq:     req.SendReq,
		StudentList: req.StudentList,
		UpdateTime:  utils.CurrentTime(),
		CreateTime:  utils.CurrentTime(),
	}
	return mysql.NewExamMysql().SendExamStudent(se)
}

// SendClass ,按班级发放试卷
//Param：
// ClassID：班级编号
//Return value：
// 错误信息
func (ec *ExamService) SendClass(req model.SendClassReq) error {
	ce := model.ClassExam{
		SendReq:    req.SendReq,
		ClassList:  req.ClassList,
		UpdateTime: utils.CurrentTime(),
		CreateTime: utils.CurrentTime(),
	}
	return mysql.NewExamMysql().SendExamClass(ce)
}
