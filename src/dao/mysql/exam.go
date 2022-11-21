package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"personality-teaching/src/model"
	"strings"
)

type examFunc interface {
	Insert(exam model.Exam) error
	UpdateExam(exam model.Exam) error
	DeleteExam(examID string) error
	QueryExam(examID string) (model.ExamDetailResp, error)
	QueryExamList(teacherID string, offset int, pageSize int) (model.ExamListResp, error)
	SendExamStudent(req model.StudentExam) error
	SendExamClass(req model.ClassExam) error
	Query(text string, teacherID string) (model.ExamListResp, error)
}

type ExamMySQL struct{}

// Query 模糊查询
func (e ExamMySQL) Query(text string, teacherID string) (model.ExamListResp, error) {
	var exams []model.ExamResp
	examName := "%" + text + "%"
	// 查询试卷列表
	if err := db.Raw("select `exam_id`,`exam_name`,`comment`,`update_time` "+
		"from `t_exam` "+
		"where `create_teacher_id`=? and `exam_name` like ? "+
		"order by `create_time` "+
		"desc ",
		teacherID, examName).Scan(&exams).Error; err != nil {
		return model.ExamListResp{}, err
	}
	var examList model.ExamListResp
	examList.Total = len(exams)
	examList.ExamList = exams
	return examList, nil
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
func (e ExamMySQL) QueryExamList(teacherID string, offset int, pageSize int) (model.ExamListResp, error) {
	var exams []model.ExamResp
	// 查询试卷列表
	if err := db.Raw("select `exam_id`,`exam_name`,`comment`,`update_time` "+
		"from `t_exam` "+
		"where `create_teacher_id`=? "+
		"order by `create_time` "+
		"desc "+
		"limit ?,?",
		teacherID, offset, pageSize).Scan(&exams).Error; err != nil {
		return model.ExamListResp{}, err
	}

	var total int
	// 查询试卷总数
	if err := db.Raw("select count(*) from `t_exam` where `create_teacher_id`=?", teacherID).Scan(&total).Error; err != nil {
		return model.ExamListResp{}, err
	}
	var examList model.ExamListResp
	examList.Total = total

	examList.ExamList = exams

	return examList, nil
}

// SendExamStudent ,插入学生试卷
func (e ExamMySQL) SendExamStudent(req model.StudentExam) error {
	var sql string
	sql = "insert into `t_student_exam`(`exam_id`,`student_id`,`comment`,`start_time`,`end_time`,`update_time`,`create_time`) values"
	for k, _ := range req.StudentList {
		var temp string
		if k == len(req.StudentList)-1 {
			temp = fmt.Sprintf("%s'%s','%s','%s','%s','%s','%s','%s'%s;", "(", req.ExamID, req.StudentList[k].StudentID, req.Comment, req.StartTime, req.EndTime, req.UpdateTime, req.CreateTime, ")")
		} else {
			temp = fmt.Sprintf("%s'%s','%s','%s','%s','%s','%s','%s'%s,", "(", req.ExamID, req.StudentList[k].StudentID, req.Comment, req.StartTime, req.EndTime, req.UpdateTime, req.CreateTime, ")")
		}
		var build strings.Builder
		build.WriteString(sql)
		build.WriteString(temp)
		sql = build.String()
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(sql).Error; err != nil {
			return err
		}
		return nil
	})
}

// SendExamClass ,按班级插入试卷
func (e ExamMySQL) SendExamClass(ce model.ClassExam) error {
	pre := "INSERT `t_student_exam`(`exam_id`,`t_student_exam`.`student_id`,`t_student_exam`.`comment`,`start_time`,`end_time`,`update_time`,`create_time`) SELECT ?,`student_id`,?,?,?,?,? FROM `t_student`WHERE `class_id` in("
	var temp string

	for k, _ := range ce.ClassList {
		if k == len(ce.ClassList)-1 {
			temp = fmt.Sprintf("'%s'", ce.ClassList[k].ClassID)
		} else {
			temp = fmt.Sprintf("'%s',", ce.ClassList[k].ClassID)
		}
		var build strings.Builder
		build.WriteString(pre)
		build.WriteString(temp)
		pre = build.String()
	}
	sql := fmt.Sprintf("%s%s", pre, ")")
	fmt.Println(sql)

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(sql, ce.ExamID, ce.Comment, ce.StartTime, ce.EndTime, ce.UpdateTime, ce.CreateTime).Error; err != nil {
			return err
		}
		return nil
	})
}

var _ examFunc = &ExamMySQL{}

func NewExamMysql() *ExamMySQL {
	return &ExamMySQL{}
}
