package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {

	//初始化logger
	initialize.InitLogger()

	//初始化 配置文件
	initialize.InitConfig()

	//初始化路由
	Router := initialize.Routers()

	//初始化翻译
	err := initialize.InitValidator("zh")
	if err != nil {
		zap.S().Infof("初始化翻译失败")
		return
	}

	port := global.ServerConfig.Port

	//S和L函数很有用
	zap.S().Infof("启动服务器，端口：%d", port)
	err = Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
