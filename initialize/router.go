package initialize

import (
	"github.com/gin-gonic/gin"
	"groqai2api/middlewares"
	"groqai2api/router"
)

func InitRouter() *gin.Engine {
	Router := gin.Default()

	Router.Use(middlewares.Cors)

	Router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "https://github.com/Star-Studio-Develop/Groq2API")
	})

	Router.GET("/ping", func(c *gin.Context) {
		c.String(200, "(｡•ˇ‸ˇ•｡)哼！都怪你\n(`ȏ´) 也不哄哄人家\n(〃′o`)人家超想哭的，捶你胸口，大坏蛋！！！\n(｡•ˇ‸ˇ•｡)哼！都怪你\n(`ȏ´)也不哄哄人家\n(〃′o`)人家超想哭的，捶你胸口，大坏蛋！！！\n(￣^￣)ゞ咩QAQ捶你胸口你好讨厌！\n(￣^￣)ゞ咩QAQ捶你胸口你好讨厌！\n(=ﾟωﾟ)ﾉ要抱抱嘤嘤嘤哼，要抱抱嘤嘤嘤哼，人家拿小拳拳捶你胸口！！！\n(=ﾟωﾟ)ﾉ要抱抱嘤嘤嘤哼，要抱抱嘤嘤嘤哼，人家拿小拳拳捶你胸口！！！\n(｡• ︿•̀｡)大坏蛋，打死你(つд⊂)")
	})
	v1Group := Router.Group("/v1/")
	router.InitRouter(v1Group)

	return Router
}
