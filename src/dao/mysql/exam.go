package mysql

import (
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type examFunc interface {
	Insert(exam model.Exam) error
	UpdateExam(exam model.Exam) error
	DeleteExam(examID string) error
	QueryExam(examID string) (model.ExamDetailResp, error)
	QueryExamList(teacherID string, page int, pageSize int) ([]model.ExamResp, error)
	SendExamStudent(req model.StudentExam) error
	SendExamClass(req model.ClassExam) error
}

type ExamMySQL struct{}

var _ examFunc = &ExamMySQL{}

func NewExamMysql() *ExamMySQL {
	return &ExamMySQL{}
}

// Insert 插入试卷
func (e ExamMySQL) Insert(exam model.Exam) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("insert into `t_exam`(`exam_id`,`exam_name`,`questions`,`comment`,`create_teacher_id`,`update_time`,`create_time`) values(?,?,?,?,?,?,?)",
			exam.ExamID, exam.ExamName, exam.Questions, exam.Comment, exam.CreateTeacherID, exam.UpdateTime, exam.CreateTime).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateExam 更新试卷
func (e ExamMySQL) UpdateExam(exam model.Exam) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("update `t_exam` set `exam_name`=?,`questions`=?,`comment`=?,`update_time`=? where `exam_id`=?",
			exam.ExamName, exam.Questions, exam.Comment, exam.UpdateTime, exam.ExamID).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteExam ,删除试卷
func (e ExamMySQL) DeleteExam(examID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("delete from `t_exam` where `exam_id`=?", examID).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryExam ,获取试卷详细消息
func (e ExamMySQL) QueryExam(examID string) (model.ExamDetailResp, error) {
	var exam model.ExamDetailResp
	if err := db.Raw("select `exam_id`,`exam_name`,`questions`,`comment`,`update_time` from `t_exam` where `exam_id`=?",
		examID).Scan(&exam).Error; err != nil {
		return model.ExamDetailResp{}, err
	}
	return exam, nil
}

// QueryExamList ,获取老师试卷列表
// Param:
// teacherID: 老师编号
// offset: 开始的序号（最小为0开始）
// pageSize: 每页的大小
func (e ExamMySQL) QueryExamList(teacherID string, offset int, pageSize int) ([]model.ExamResp, error) {
	var exams []model.ExamResp
	if err := db.Raw("select `exam_id`,`exam_name`,`comment`,`update_time` "+
		"from `t_exam` "+
		"where `create_teacher_id`=? "+
		"order by `create_time` "+
		"desc "+
		"limit ?,?",
		teacherID, offset, pageSize).Scan(&exams).Error; err != nil {
		return nil, err
	}
	return exams, nil
}

// SendExamStudent ,插入学生试卷
func (e ExamMySQL) SendExamStudent(req model.StudentExam) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("insert into `t_student_exam`(`exam_id`,`student_id`,`comment`,`start_time`,`end_time`,`update_time`,`create_time`) values(?,?,?,?,?,?,?)",
			req.ExamID, req.StudentID, req.Comment, req.StartTime, req.EndTime, req.UpdateTime, req.CreateTime).Error; err != nil {
			return err
		}
		return nil
	})
}

// SendExamClass ,按班级插入试卷
func (e ExamMySQL) SendExamClass(ce model.ClassExam) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("INSERT `t_student_exam`(`exam_id`,`t_student_exam`.`student_id`,`t_student_exam`.`comment`,`start_time`,`end_time`,`update_time`,`create_time`) "+
			"SELECT ?,`student_id`,?,?,?,?,? "+
			"FROM `t_student` "+
			"WHERE `class_id`=? ",
			ce.ExamID, ce.Comment, ce.StartTime, ce.EndTime, ce.UpdateTime, ce.CreateTime, ce.ClassID).Error; err != nil {
			return err
		}
		return nil
	})
}
