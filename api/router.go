package api

import (
	"github.com/gin-gonic/gin"
	"groqai2api/initialize"
	"net/http"
)

var router *gin.Engine

func init() {
	// 初始化配置
	initialize.InitConfig()
	// 初始化缓存
	initialize.InitCache()
	// 初始化代理
	initialize.InitProxy()
	// 初始化账号
	initialize.InitAuth()
	// 初始化gin
	router = initialize.InitRouter()
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
