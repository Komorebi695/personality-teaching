package main

import (
	"github.com/gin-gonic/gin"
	"personality-teaching/src/configs"
	"personality-teaching/src/dao"
)

func main() {
	// 创建数据库连接
	if err := dao.InitMysql(); err != nil {
		panic(err.Error())
	}

	r := gin.Default()

	// 监听端口
	config := configs.InitConfig()
	addr := ":" + config.Port
	_ = r.Run(addr)
}
