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
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE `t_student_exam` SET `answers`=?,`status`=?,`score`=?,`update_time`=? WHERE `exam_id`=? and `student_id`=?",
			exams.Answers, exams.Status, exams.Score, exams.UpdateTime, exams.ExamID, exams.StudentID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r ReviewMySQL) QueryStudent(examID string, studentID string) (model.StudentExams, error) {
	var studentExam model.StudentExams
	if err := db.Raw("SELECT se.`exam_id`,`student_id`,`exam_name`,`answers`,`score`,se.`comment` FROM `t_student_exam` se LEFT JOIN `t_exam` e ON se.`exam_id`=e.`exam_id` WHERE se.`exam_id`=? AND `student_id`=? limit 0,1",
		examID, studentID).Scan(&studentExam).Error; err != nil {
		return model.StudentExams{}, err
	}
	return studentExam, nil
}

func (r ReviewMySQL) QueryStudentList(classID string, examID string) ([]model.ReviewStudent, error) {
	var studentList []model.ReviewStudent
	if err := db.Raw("SELECT s.`name`,`score`,`status`,`update_time` FROM `t_student_exam` se LEFT JOIN `t_student` s ON se.`student_id`=s.`student_id` WHERE `class_id`=? AND `exam_id`=? ORDER BY `status`;",
		classID, examID).Scan(&studentList).Error; err != nil {
		return nil, err
	}
	return studentList, nil
}

func (r ReviewMySQL) QueryClass(examID string) ([]model.ReviewClass, error) {
	var classList []model.ReviewClass
	if err := db.Raw("SELECT s.`class_id`,c.`name`,c.`college`,c.`major`,COUNT(*) AS `count` "+
		"\nFROM `t_student_exam` se "+
		"\nLEFT JOIN `t_student` s "+
		"\nON se.`student_id`=s.`student_id` "+
		"\nLEFT JOIN `t_class` c "+
		"\nON s.`class_id`=c.`class_id` "+
		"\nWHERE `status`='1' AND `exam_id`=? "+
		"\nGROUP BY s.`class_id`;", examID).Scan(&classList).Error; err != nil {
		return nil, err
	}
	return classList, nil
}

var _ reviewFunc = &ReviewMySQL{}

func NewReviewMysql() *ReviewMySQL {
	return &ReviewMySQL{}
}
