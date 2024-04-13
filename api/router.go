package api

import (
	"net/http"

	"github.com/Star-Studio-Develop/Groq2API/initialize"
	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	// 初始化全局router，这假设initialize包中有一个返回*mux.Router的SetupRouter函数
	router = initialize.SetupRouter()
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
