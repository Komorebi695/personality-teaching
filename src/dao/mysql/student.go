package mysql

import "personality-teaching/src/model"

type studentFunc interface {
	InsertStudent(student model.Student) error
	UpdateClassID(studentID, classID string) error
	QueryStudent(studentID string) (model.Student, error)
	QueryStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, int, error)
	CheckStudentClass(studentID string, classID string) (bool, error)
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

func (s *StudentMySQL) QueryStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, int, error) {
	var students []model.ClassStudentListResp
	var total int
	offset := (req.PageNum - 1) * req.PageSize
	count := req.PageSize
	err := db.Raw("select `student_id`,`name`,`college`,`major`,`phone_number` from `t_student` where `class_id` = ? limit ?,?", req.ClassID, offset, count).Scan(&students).Error
	if err != nil {
		return []model.ClassStudentListResp{}, 0, err
	}
	err = db.Raw("select count(`id`) from `t_student` where `class_id` = ?", req.ClassID).Scan(&total).Error
	if err != nil {
		return []model.ClassStudentListResp{}, 0, err
	}
	return students, total, nil
}

// CheckStudentClass 检查学生是否在班级中
func (s *StudentMySQL) CheckStudentClass(studentID string, classID string) (bool, error) {
	id := -1
	err := db.Raw("select `id` from `t_student` where `student_id` = ? and `class_id` = ?", studentID, classID).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id != -1, nil
}
