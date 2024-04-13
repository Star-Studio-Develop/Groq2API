package router

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
)

// SetupRouter 初始化和返回*mux.Router
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// 创建缓存实例
	c := cache.New(4*time.Minute, 1*time.Minute)

	// 设置根路由重定向
	r.HandleFunc("/", rootHandler).Methods("GET")

	// 设置API端点
	r.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		chatCompletionsHandler(w, r, c)
	}).Methods("POST")

	return r
}

// StartServer 启动HTTP服务器
func StartServer() {
	r := SetupRouter()

	// 获取环境变量中的端口或使用默认值
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Msgf("Server is listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
