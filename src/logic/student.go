package logic

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmdrws/golang_common/lib"
	"personality-teaching/src/dao"
	"personality-teaching/src/middleware"
	"personality-teaching/src/model"
)

type StudentService struct {
}

func (ss *StudentService) StudentList(c *gin.Context, params *model.StudentListInput) (*model.StudentListOutput, error) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return nil, err
	}
	//从db中分页读取基本信息
	questionInfo := &dao.TStudent{}
	list, total, err := questionInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return nil, err
	}

	//格式化输出信息
	var outList []model.StudentListItemOutput
	for _, listItem := range list {
		////获取学生班级
		//classItem := &dao.TClass{ClassId: listItem.ClassId}
		//classItem, err = classItem.Find(c, tx, classItem)
		//if err != nil {
		//	middleware.ResponseError(c, 2003, err)
		//	return
		//}
		outItem := model.StudentListItemOutput{
			StudentID:   listItem.StudentId,
			StudentName: listItem.Name,
			College:     listItem.College,
			Major:       listItem.Major,
			//Class:       classItem.Name,
			PhoneNumber: listItem.PhoneNumber,
		}
		outList = append(outList, outItem)
	}

	out := &model.StudentListOutput{
		Total: total,
		List:  outList,
	}
	return out, nil
}

func (ss *StudentService) StudentDelete(c *gin.Context, params *model.StudentDeleteInput) error {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return nil
	}

	//读取基本信息
	studentInfo := &dao.TStudent{StudentId: params.StudentId}
	studentInfo, err = studentInfo.Find(c, tx, studentInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New("该学生不存在"))
		return err
	}
	studentInfo.IsDelete = 1
	if err := studentInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return err
	}
	return nil
}

func (ss *StudentService) StudentAdd(c *gin.Context, params *model.StudentAddInput) error {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return nil
	}

	//读取基本信息
	studentInfo := &dao.TStudent{StudentId: params.StudentId}
	studentInfo, err = studentInfo.Find(c, tx, studentInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New("该学生不存在"))
		return err
	}
	tx = tx.Begin()

	//包装学生信息
	studentModel := &dao.TStudent{
		StudentId:   params.StudentId,
		Name:        params.StudentName,
		Password:    params.Password,
		College:     params.College,
		Major:       params.Major,
		ClassId:     params.ClassId,
		PhoneNumber: params.PhoneNumber,
		IsDelete:    0,
	}

	if err := studentModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return err
	}
	tx.Commit()
	return nil
}
