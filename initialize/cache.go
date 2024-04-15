package initialize

import (
	"github.com/patrickmn/go-cache"
	"groqai2api/global"
	"time"
)

func InitCache() {
	global.Cache = cache.New(5*time.Minute, 10*time.Minute)
}
