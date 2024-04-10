package stream

import (
	"Groq2API/initialize/model"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FetchStream 异步获取数据流
func FetchStream(jwt string, orgID string, messages []model.Message, modelType string) ([]string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	payload := map[string]interface{}{
		"model":       modelType,
		"messages":    messages,
		"temperature": 0.2,
		"max_tokens":  2048,
		"top_p":       0.8,
		"stream":      true,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("groq-organization", orgID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("Stream started...")
	fmt.Println("Response status:", resp.Status)
	var response []string
	reader := bufio.NewReader(resp.Body)
	for {
		var line []byte
		line, err = reader.ReadBytes('\n')
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			fmt.Println("Error while reading response:", err)
		}
		// log.Println(string(line))
		response = append(response, string(line))
	}
	return response, err
}
