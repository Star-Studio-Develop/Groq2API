package initialize

import (
	"groqai2api/global"
	"os"

	"github.com/joho/godotenv"
)

func InitConfig() {
	_ = godotenv.Load(".env")
	global.Host = os.Getenv("SERVER_HOST")
	if global.Host == "" {
		global.Host = "0.0.0.0"
	}
	global.Port = os.Getenv("SERVER_PORT")
	if global.Port == "" {
		global.Port = os.Getenv("PORT")
		if global.Port == "" {
			global.Port = "8080"
		}
	}

	global.ChinaPrompt = os.Getenv("CHINA_PROMPT")
	global.Authorization = os.Getenv("Authorization")
}
