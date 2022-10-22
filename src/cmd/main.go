package main

import (
	"github.com/jmdrws/golang_common/lib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化配置文件
	lib.InitModule("./src/configs/", []string{"base", "mysql", "redis"})
	defer lib.Destroy()

	r := InitRouter()
	go func() {
		// 监听端口
		addr := lib.GetStringConf("base.http.addr")
		log.Printf(" [INFO] HttpServerRun:%s\n", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf(" [ERROR] ServerRun:%s err:%v\n", addr, err)
		}
	}()

	//quit持续监听信号（syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
