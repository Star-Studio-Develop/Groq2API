package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// TokenResponse 定义了JWT响应的结构
type TokenResponse struct {
	Data struct {
		SessionJWT string `json:"session_jwt"`
	} `json:"data"`
}

// FetchJWT 使用HTTP请求获取JWT
func FetchJWT(refreshToken string) (string, error) {
	url := "https://web.stytch.com/sdk/v1/sessions/authenticate"
	body := bytes.NewBuffer(nil)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+refreshToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.Data.SessionJWT, nil
}
