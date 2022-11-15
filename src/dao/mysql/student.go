package mysql

import "personality-teaching/src/model"

type studentFunc interface {
	InsertStudent(student model.Student) error
	UpdateClassID(studentID, classID string) error
	QueryStudent(studentID string) (model.Student, error)
	QueryStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, error)
}

var _ studentFunc = &StudentMySQL{}

type StudentMySQL struct{}

func NewStudentMySQL() *StudentMySQL {
	return &StudentMySQL{}
}

func (s *StudentMySQL) InsertStudent(stu model.Student) error {
	return db.Create(&stu).Error
}

func (s *StudentMySQL) UpdateClassID(studentID, classID string) error {
	return db.Exec("update `t_student` set `class_id` = ? where student_id = ?", classID, studentID).Error
}

func (s *StudentMySQL) QueryStudent(studentID string) (model.Student, error) {
	var m model.Student
	err := db.Raw("select `student_id`,`name`,`password`,`college`,`major`,`class_id`,`phone_number` from `t_student` where `student_id` = ?", studentID).Scan(&m).Error
	if err != nil {
		return model.Student{}, err
	}
	return m, nil
}

func (s *StudentMySQL) QueryStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, error) {
	var students []model.ClassStudentListResp
	offset := (req.PageNum - 1) * req.PageSize
	count := req.PageSize
	err := db.Raw("select `student_id`,`name`,`college`,`major`,`phone_number` from `t_student` where `class_id` = ? limit ?,?", req.ClassID, offset, count).Scan(&students).Error
	if err != nil {
		return []model.ClassStudentListResp{}, err
	}
	return students, nil
}
