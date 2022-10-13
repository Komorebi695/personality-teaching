package main

import (
	"github.com/gin-gonic/gin"
	"personality-teaching/src/configs"
	"personality-teaching/src/dao"
)

func main() {
	// 初始化配置文件
	config := configs.InitConfig()

	// 创建数据库连接
	if err := dao.InitMysql(config); err != nil {
		panic(err.Error())
	}

	r := gin.Default()

	// 监听端口
	addr := ":" + config.Port
	_ = r.Run(addr)
}
