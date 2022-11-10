package redis

import (
	"time"

	"github.com/go-redis/redis"
)

//Key: session ID，UUID随机串，前端cookie保存。
//Value: json结构体序列化存储，结构体中包含如下字段
//◦ UserID: ⽤⼾ID
//◦ RoleType: 学⽣/教师
//◦ CreateTime: 创建时间(时间戳)

const expireTime = time.Hour * 24 * 7

func SetSessionNX(key string, value interface{}) error {
	return redisDB.SetNX(key, value, expireTime).Err()
}

// GetSessionValue 存在key返回value值，不存在返回空字符串
func GetSessionValue(sessionKey string) (string, error) {
	val, err := redisDB.Get(sessionKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func ResetExpireTime(key string) error {
	return redisDB.Expire(key, expireTime).Err()
}
