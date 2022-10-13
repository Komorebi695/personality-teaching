package main

import (
	"fmt"
	"personality-teaching/src/configs"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/dao/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置文件
	config := configs.InitConfig()
	if config == nil {
		panic("配置文件加载错误!")
	}

	// 创建数据库连接
	if err := mysql.InitMysql(config); err != nil {
		panic("MySQL init error: " + err.Error())
	}

	if err := redis.InitRedis(config.Redis); err != nil {
		panic("Redis init error: " + err.Error())
	}
	r := gin.Default()

	// 监听端口
	addr := fmt.Sprintf(":%s", config.Port)
	_ = r.Run(addr)
}
