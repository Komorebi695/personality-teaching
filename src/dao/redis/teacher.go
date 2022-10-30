package redis

import (
	"fmt"
	"time"
)

const teacherIDPre = "teacher"

// session_key在Redis中value为string 类型：
//key --> teacher:{hashed Teacher ID}   value --> {teacherId}
func getKey(teacherID string) string {
	return fmt.Sprintf("%s:%s", teacherIDPre, teacherID)
}

func SetSessionNX(teacherID string, value interface{}) error {
	return redisDB.SetNX(getKey(teacherID), value, time.Hour*240).Err()
}
