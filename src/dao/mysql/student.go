package mysql

import (
	"fmt"
	"personality-teaching/src/model"
)

type studentFunc interface {
	InsertStudent(student model.Student) error
	UpdateClassID(studentID, classID string) error
	QueryStudent(studentID string) (model.Student, error)
	QueryStudentsInClass(req model.ClassStudentListReq) ([]model.ClassStudentListResp, int, error)
	QueryStudentsNotInClass(req model.ClassStudentListReq, content string) ([]model.ClassStudentListResp, int, error)
	CheckStudentClass(studentID string, classID string) (bool, error)
	UpdatePassWord(studentID string, newPwd string) error
	QueryAllByName(name string) (model.Student, error)
	QueryStudentLike(searchText string) ([]model.ClassStudentListResp, error)
	DeleteStudent(studentID string) error
	UpdateStudent(cs model.CreateStudentResp) error
}

var _ studentFunc = &StudentMySQL{}

type StudentMySQL struct{}

func NewStudentMySQL() *StudentMySQL {
	return &StudentMySQL{}
}

func (s *StudentMySQL) InsertStudent(stu model.Student) error {
	return Db.Create(&stu).Error
}

func (s *StudentMySQL) DeleteStudent(studentID string) error {
	return Db.Exec("delete from `t_student` where `student_id`=?", studentID).Error
}

func (s *StudentMySQL) UpdateClassID(studentID, classID string) error {
	return Db.Exec("update `t_student` set `class_id` = ? where student_id = ?", classID, studentID).Error
}

func (s *StudentMySQL) QueryStudent(studentID string) (model.Student, error) {
	var m model.Student
	err := Db.Raw("select `student_id`,`name`,`student_no`,`password`,`college`,`major`,`class_id`,`phone_number` from `t_student` where `student_id` = ?", studentID).Scan(&m).Error
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
	err := Db.Raw("select `student_id`,`name`,`student_no`,`college`,`major`,`phone_number` from `t_student` where `class_id` = ? limit ?,?", req.ClassID, offset, count).Scan(&students).Error
	if err != nil {
		return []model.ClassStudentListResp{}, 0, err
	}
	err = Db.Raw("select count(`id`) from `t_student` where `class_id` = ?", req.ClassID).Scan(&total).Error
	if err != nil {
		return []model.ClassStudentListResp{}, 0, err
	}
	return students, total, nil
}

func (s *StudentMySQL) QueryStudentsNotInClass(req model.ClassStudentListReq, content string) ([]model.ClassStudentListResp, int, error) {
	var students []model.ClassStudentListResp
	var total int
	offset := (req.PageNum - 1) * req.PageSize
	count := req.PageSize
	content = fmt.Sprintf("%%%s%%", content)
	err := Db.Raw("SELECT `student_id`,`name`,`student_no`,`college`,`major`,`phone_number` FROM `t_student` WHERE `class_id` = ? AND `name` LIKE ?  LIMIT ?,?", req.ClassID, content, offset, count).Scan(&students).Error
	if err != nil {
		return []model.ClassStudentListResp{}, 0, err
	}
	err = Db.Raw("select count(`id`) from `t_student` where `class_id` = ? AND `name` LIKE ?;", req.ClassID, content).Scan(&total).Error
	if err != nil {
		return []model.ClassStudentListResp{}, 0, err
	}
	return students, total, nil
}

// CheckStudentClass 检查学生是否在班级中
func (s *StudentMySQL) CheckStudentClass(studentID string, classID string) (bool, error) {
	id := -1
	err := Db.Raw("select `id` from `t_student` where `student_id` = ? and `class_id` = ?", studentID, classID).Scan(&id).Error
	if err != nil {
		return false, err
	}
	return id != -1, nil
}

func (s *StudentMySQL) UpdatePassWord(studentID string, newPwd string) error {
	return Db.Exec("update `t_student` set `password` = ? where `student_id` = ?", newPwd, studentID).Error
}

func (s *StudentMySQL) QueryAllByName(name string) (model.Student, error) {
	var m model.Student
	err := Db.Raw("select `student_id`,`name`,`student_no`,`password`,`college`,`major`,`class_id`,`phone_number` from `t_student` where `name` = ? limit 1", name).Scan(&m).Error
	if err != nil {
		return model.Student{}, err
	}
	return m, nil
}

func (s *StudentMySQL) QueryStudentLike(searchText string) ([]model.ClassStudentListResp, error) {
	var m []model.ClassStudentListResp
	searchText = fmt.Sprintf("%%%s%%", searchText)
	err := Db.Raw("select `student_id`,`name`,`student_no`,`college`,`major`,`phone_number` from `t_student` where `class_id`='0' and `name` like ?;", searchText).Scan(&m).Error
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		err := Db.Raw("select `student_id`,`name`,`student_no`,`college`,`major`,`phone_number` from `t_student` where `class_id`='0' and `student_no` like ?;", searchText).Scan(&m).Error
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *StudentMySQL) UpdateStudent(cs model.CreateStudentResp) error {
	return Db.Exec("update `t_student` set `name`=?,`student_no`=?,`college`=?,`major`=?,`phone_number`=? where `student_id`=?;", cs.Name, cs.StudentNo, cs.College, cs.Major, cs.PhoneNumber, cs.StudentID).Error
}

func (s *StudentMySQL) QueryStudentInStuQu(studentID string) (stu []model.StudentQuestion, err error) {
	err = Db.Raw("select knp_id, allscore, score,answer from t_student_question q LEFT JOIN t_knowledge_point_question knp ON q.question_id = knp.question_id where student_id = ?", studentID).Scan(&stu).Error
	if err != nil {
		return nil, err
	}
	return stu, nil
}

func (s *StudentMySQL) QueryAllKnp() (knp []model.Studentknp, err error) {
	err = Db.Raw("SELECT knp_id,name,parent_knp_id,level FROM t_knowledge_point").Scan(&knp).Error
	if err != nil {
		return nil, err
	}
	return knp, nil
}

func (s *StudentMySQL) QueryteacherClass(teacherID string) (class []model.StudentClass, err error) {
	err = Db.Raw("select s.class_id, s.student_id,c.name from t_student s LEFT JOIN t_teacher_class t ON t.class_id = s.class_id LEFT JOIN t_class c ON c.class_id = s.class_id where teacher_id = ? and is_valid = 1", teacherID).Scan(&class).Error
	if err != nil {
		return nil, err
	}
	return class, nil
}

func (s *StudentMySQL) Queryteacherstudent(ClassID string) (stu []model.StudentIDandName, err error) {
	err = Db.Raw("select student_id,`name` from t_student WHERE class_id = ?", ClassID).Scan(&stu).Error
	if err != nil {
		return nil, err
	}
	return stu, nil
}
