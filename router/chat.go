package router

import (
	"bufio"
	"encoding/json"
	"groqai2api/global"
	groqHttp "groqai2api/pkg/groq"
	"io/ioutil"
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
	var apiReq groq.APIRequest
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
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
			// Handle session token
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
			// Handle custom token
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

	// Insert Chinese prompt if needed
	if global.ChinaPrompt == "true" {
		prompt := groq.APIMessage{
			Content: "使用中文回答，输出时不要带英文",
			Role:    "system",
		}
		apiReq.Messages = append([]groq.APIMessage{prompt}, apiReq.Messages...)
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

	apiKey, _ := global.Cache.Get(account.Organization)
	response, err := groqHttp.ChatCompletions(client, apiReq, apiKey.(string), account.Organization, "")
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	defer response.Body.Close()

	// Check if the client supports streaming
	if strings.Contains(c.GetHeader("Accept"), "text/event-stream") {
		// Stream the response
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		reader := bufio.NewReader(response.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				break
			}
			c.Writer.Write(line)
			c.Writer.Flush()
		}
	} else {
		// Non-streaming response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(body)
	}
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
