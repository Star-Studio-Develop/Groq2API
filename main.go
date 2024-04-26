package main

import (
	"groqai2api/global"
	"groqai2api/initialize"
)

func main() {
	// 初始化配置
	initialize.InitConfig()
	// 初始化缓存
	initialize.InitCache()
	// 初始化代理
	initialize.InitProxy()
	// 初始化账号
	initialize.InitAuth()
	// 初始化路由
	Router := initialize.InitRouter()
	if err := Router.Run(global.Host + ":" + global.Port); err != nil {
		panic(err)
	}
}
