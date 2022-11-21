package logic

import (
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"
)

type ReviewService struct{}

func NewReviewService() *ReviewService {
	return &ReviewService{}
}

func (rs *ReviewService) UpdateReview(review model.ReviewUpdateReq) error {
	var reviewStudent = model.ReviewUpdate{
		ReviewUpdateReq: review,
		UpdateTime:      utils.CurrentTime(),
	}
	return mysql.NewReviewMysql().UpdateReview(reviewStudent)
}

func (rs *ReviewService) QueryReviewClass(examID string) ([]model.ReviewClass, error) {
	return mysql.NewReviewMysql().QueryClass(examID)
}

func (rs *ReviewService) QueryReviewStudent(classID string, examID string) ([]model.ReviewStudent, error) {
	return mysql.NewReviewMysql().QueryStudentList(classID, examID)
}

func (rs *ReviewService) QueryStudentAnswer(examID string, studentId string) (model.StudentExams, error) {
	return mysql.NewReviewMysql().QueryStudent(examID, studentId)
}
