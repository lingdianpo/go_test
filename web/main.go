package main

import (
	"fmt"
	"go.uber.org/zap"
	"test/web/initialize"
)

func main() {
	port := 8021
	//1. 初始化logger
	initialize.LoggerInit()
	//2. 初始化routers
	Router := initialize.RouterInit()
	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分几倍的，debug，info，warn，error，fetal
		3.S函数和L函数很有用,提供了一个全局的安全访问logger的途径
	*/
	zap.S().Info("启动服务器，端口 %d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Error("启动失败", zap.Any("err", err))
	}

}
