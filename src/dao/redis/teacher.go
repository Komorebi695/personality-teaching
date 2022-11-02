package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const teacherIDPre = "teacher"

// session_key在Redis中value为string 类型：
//key --> teacher:{hashed Teacher ID}   value --> {teacherId}
func getTeacherKey(hashedTeacherID string) string {
	return fmt.Sprintf("%s:%s", teacherIDPre, hashedTeacherID)
}

func SetSessionNX(teacherID string, value interface{}) error {
	return redisDB.SetNX(getTeacherKey(teacherID), value, time.Hour*240).Err()
}

// GetTeacherIX 存在key返回value值，不存在返回空字符串
func GetTeacherIX(sessionKey string) (string, error) {
	val, err := redisDB.Get(getTeacherKey(sessionKey)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}
