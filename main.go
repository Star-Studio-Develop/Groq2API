package main

import (
	"groqai2api/global"
	"groqai2api/initialize"
)

func main() {
	//jar := tls_client.NewCookieJar()
	//options := []tls_client.HttpClientOption{
	//	tls_client.WithTimeoutSeconds(30),
	//	tls_client.WithClientProfile(profiles.Okhttp4Android13),
	//	tls_client.WithNotFollowRedirects(),
	//	tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	//}

	//client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//req, err := http.NewRequest(http.MethodGet, "https://api.groq.com/platform/v1/user/profile", nil)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//

	//println(client)
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//if resp.StatusCode != 200 {
	//	println(resp.StatusCode)
	//	log.Fatal("not valid response")
	//}
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
