package stream

import (
	"bytes"
	"encoding/json"
	"github.com/Star-Studio-Develop/Groq2API/initialize/model"
	"github.com/rs/zerolog/log"
	"net/http"
)

// FetchStream 异步获取数据流
func FetchStream(jwt string, orgID string, messages []model.Message, modelType string, maxTokens int64) (*http.Response, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	payload := map[string]interface{}{
		"model":       modelType,
		"messages":    messages,
		"temperature": 0.2,
		"max_tokens":  maxTokens,
		"top_p":       0.8,
		"stream":      true,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Error().Err(err).Msg("Error creating request: ")
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("groq-organization", orgID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending request: ")
		return nil, err
	}
	//_ = resp.Body.Close()
	_ = req.Body.Close()

	return resp, nil
}
