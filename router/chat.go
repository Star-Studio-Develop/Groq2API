package router

import (
	"encoding/json"
	"groqai2api/global"
	groqHttp "groqai2api/pkg/groq"
	"log/slog"
	"net/http"
	"strings"
	"time"

	tls_client "github.com/bogdanfinn/tls-client"

	"github.com/gin-gonic/gin"
	groq "github.com/learnLi/groq_client"
)

func authSessionHandler(client tls_client.HttpClient, account *groq.Account, api_key string, proxy string) error {
	organization, err := groqHttp.GerOrganizationId(client, api_key, proxy)
	if err != nil {
		slog.Error("Failed to get organization id", "err", err)
		return err
	}
	account.Organization = organization
	global.Cache.Set(organization, api_key, 3*time.Minute)
	return nil
}

func authRefreshHandler(client tls_client.HttpClient, account *groq.Account, api_key string, proxy string) error {
	token, err := groqHttp.GetSessionToken(client, api_key, "")
	if err != nil {
		slog.Error("Failed to get session token", "err", err)
		return err
	}
	organization, err := groqHttp.GerOrganizationId(client, token.Data.SessionJwt, proxy)
	if err != nil {
		slog.Error("Failed to get organization id", "err", err)
		return err
	}
	account.Organization = organization
	global.Cache.Set(organization, token.Data.SessionJwt, 3*time.Minute)
	return nil
}

func chat(c *gin.Context) {
	var api_req groq.APIRequest
	if err := c.ShouldBindJSON(&api_req); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	client := groqHttp.NewBasicClient()
	proxyUrl := global.ProxyPool.GetProxyIP()
	if proxyUrl != "" {
		client.SetProxy(proxyUrl)
	}

	authorization := c.Request.Header.Get("Authorization")
	account := global.AccountPool.Get()
	if authorization != "" {
		customToken := strings.Replace(authorization, "Bearer ", "", 1)
		if customToken != "" {
			// 说明传递的是session_token
			if strings.HasPrefix(customToken, "eyJhbGciOiJSUzI1NiI") {
				account = groq.NewAccount("", "")
				err := authSessionHandler(client, account, customToken, "")
				if err != nil {
					slog.Error("session_token is invalid", err)
					c.JSON(400, gin.H{"error": err.Error()})
					c.Abort()
					return
				}
			}
			if len(customToken) == 44 {
				account = groq.NewAccount(customToken, "")
				err := authRefreshHandler(client, account, customToken, "")
				if err != nil {
					slog.Error("customToken is invalid", err)
					c.JSON(400, gin.H{"error": err.Error()})
					c.Abort()
					return
				}
			}
		}
	}

	// 默认插入中文prompt
	if global.ChinaPrompt == "true" {
		prompt := groq.APIMessage{
			Content: "使用中文回答，输出时不要带英文",
			Role:    "system",
		}
		api_req.Messages = append([]groq.APIMessage{prompt}, api_req.Messages...)
	}
	if _, ok := global.Cache.Get(account.Organization); !ok {
		err := authRefreshHandler(client, account, account.SessionToken, "")
		if err != nil {
			slog.Error("get refresh err", err)
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}
	api_key, _ := global.Cache.Get(account.Organization)
	response, err := groqHttp.ChatCompletions(client, api_req, api_key.(string), account.Organization, "")
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	defer response.Body.Close()
	groqHttp.NewReadWriter(c.Writer, response).StreamHandler()
}

func models(c *gin.Context) {
	client := groqHttp.NewBasicClient()
	proxyUrl := global.ProxyPool.GetProxyIP()
	if proxyUrl != "" {
		client.SetProxy(proxyUrl)
	}
	authorization := c.Request.Header.Get("Authorization")
	account := global.AccountPool.Get()
	if authorization != "" {
		customToken := strings.Replace(authorization, "Bearer ", "", 1)
		if customToken != "" {
			// 说明传递的是session_token
			if strings.HasPrefix(customToken, "eyJhbGciOiJSUzI1NiI") {
				account = groq.NewAccount("", "")
				err := authSessionHandler(client, account, customToken, "")
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					c.Abort()
					return
				}
			}
			if len(customToken) == 44 {
				account = groq.NewAccount(customToken, "")
				err := authRefreshHandler(client, account, customToken, "")
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					c.Abort()
					return
				}
			}
		}
	}

	if _, ok := global.Cache.Get(account.Organization); !ok {
		err := authRefreshHandler(client, account, account.SessionToken, "")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}
	api_key, _ := global.Cache.Get(account.Organization)
	response, err := groqHttp.GetModels(client, api_key.(string), account.Organization, "")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	var mo groq.Models

	if err := json.NewDecoder(response.Body).Decode(&mo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, mo)
}

func InitRouter(Router *gin.RouterGroup) {
	Router.GET("models", models)
	ChatRouter := Router.Group("chat")
	{
		ChatRouter.POST("/completions", chat)
	}
}
