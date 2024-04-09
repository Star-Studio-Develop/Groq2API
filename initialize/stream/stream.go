package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchStream 异步获取数据流
func FetchStream(jwt string, orgID string) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	payload := map[string]interface{}{
		"model": "mixtral-8x7b-32768",
		"messages": []map[string]string{
			{"content": "Please try to provide useful, helpful and actionable answers.", "role": "system"},
			{"content": "hi", "role": "user"},
		},
		"temperature": 0.2,
		"max_tokens":  2048,
		"top_p":       0.8,
		"stream":      true,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("groq-organization", orgID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Stream started...")
	fmt.Println("Response status:", resp.Status)
}
