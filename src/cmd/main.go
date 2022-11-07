package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"personality-teaching/src/configs"
	"personality-teaching/src/dao/mysql"
	"personality-teaching/src/dao/redis"
	"personality-teaching/src/logger"
	"syscall"
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

	//初始化日志
	logger.InitLogger()

	// 路由
	r := InitRouter()
	go func() {
		// 监听端口
		addr := fmt.Sprintf(":%s", config.Port)
		if err := r.Run(addr); err != nil {
			log.Fatalf(" [ERROR] ServerRun:%s err:%v\n", addr, err)
		}
	}()
	//quit持续监听信号（syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
