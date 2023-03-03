package mysql

import (
	"gorm.io/gorm"
	"personality-teaching/src/model"
)

type reviewFunc interface {
	QueryClass(examID string) ([]model.ReviewClass, error)
	QueryStudentList(classID string, examID string) ([]model.ReviewStudent, error)
	QueryStudent(studentID string, examID string) (model.StudentExams, error)
	UpdateReview(exams model.ReviewUpdate) error
}

type ReviewMySQL struct{}

func (r ReviewMySQL) UpdateReview(exams model.ReviewUpdate) error {
	return Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE `t_student_exam` SET `answers`=?,`detailed_score`=?,`total_score`=?,`problem_status`=?,`status`=?,`update_time`=? WHERE `exam_id`=? AND `student_id`=? AND `times`=?;",
			exams.Answers, exams.DetailedScore, exams.TotalScore, exams.ProblemStatus, exams.Status, exams.UpdateTime, exams.ExamID, exams.StudentID, exams.Times).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r ReviewMySQL) QueryStudent(examID string, studentID string) (model.StudentExams, error) {
	var studentExam model.StudentExams
	if err := Db.Raw("SELECT se.`exam_id`,`student_id`,`exam_name`,`answers`,`detailed_score`,`total_score`,`problem_status`,`times`,se.`comment` "+
		"FROM `t_student_exam` se "+
		"LEFT JOIN `t_exam` e "+
		"ON se.`exam_id`=e.`exam_id` "+
		"WHERE se.`exam_id`=? "+
		"AND `student_id`=? "+
		"ORDER BY `times` DESC "+
		"LIMIT 0,1;",
		examID, studentID).Scan(&studentExam).Error; err != nil {
		return model.StudentExams{}, err
	}
	return studentExam, nil
}

func (r ReviewMySQL) QueryStudentList(classID string, examID string) ([]model.ReviewStudent, error) {
	var studentList []model.ReviewStudent
	if err := Db.Raw("SELECT s.`student_id`,s.`name`,`total_score`,`status`,`update_time` FROM `t_student_exam` se LEFT JOIN `t_student` s ON se.`student_id`=s.`student_id` WHERE `class_id`=? AND `exam_id`=? ORDER BY `status`;",
		classID, examID).Scan(&studentList).Error; err != nil {
		return nil, err
	}
	return studentList, nil
}

func (r ReviewMySQL) QueryClass(examID string) ([]model.ReviewClass, error) {
	var classList []model.ReviewClass
	if err := Db.Raw("SELECT s.`class_id`,c.`name`,c.`college`,c.`major`,COUNT(*) AS `count` "+
		"\nFROM `t_student_exam` se "+
		"\nLEFT JOIN `t_student` s "+
		"\nON se.`student_id`=s.`student_id` "+
		"\nLEFT JOIN `t_class` c "+
		"\nON s.`class_id`=c.`class_id` "+
		"\nWHERE `exam_id`=? AND c.`name` IS NOT NULL"+
		"\nGROUP BY s.`class_id`;", examID).Scan(&classList).Error; err != nil {
		return nil, err
	}
	return classList, nil
}

var _ reviewFunc = &ReviewMySQL{}

func NewReviewMysql() *ReviewMySQL {
	return &ReviewMySQL{}
}
