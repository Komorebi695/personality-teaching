package redis

import (
	"fmt"
	"personality-teaching/src/configs"

	"github.com/go-redis/redis"
)

var redisDB *redis.Client

func InitRedis(redisConf configs.Redis) error {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConf.Addr, redisConf.Port),
		//Password: "Duxiaokun.", // redis密码，没有则留空
		DB:       0,            // 默认数据库，默认是0
	})

	//检查是否成功连接到了redis服务器
	_, err := redisDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
