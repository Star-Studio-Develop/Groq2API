package main

import (
	"bytes"
	"encoding/json"
	"github.com/Star-Studio-Develop/Groq2API/initialize/auth"
	"github.com/Star-Studio-Develop/Groq2API/initialize/model"
	"github.com/Star-Studio-Develop/Groq2API/initialize/stream"
	"github.com/Star-Studio-Develop/Groq2API/initialize/user"
	"github.com/Star-Studio-Develop/Groq2API/initialize/utils"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
	"time"
)

type ChatCompletionRequest struct {
	//RefreshToken string `json:"refresh_token"`
	Messages  []model.Message `json:"messages"`
	ModelType string          `json:"model"`
	Stream    bool            `json:"stream"`
	MaxTokens int64           `json:"max_tokens"`
}

func chatCompletionsHandler(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	var jwt string
	var err error
	var orgID string
	if r.Method == "OPTIONS" {
		utils.SetCorsHeaders(w)
		// Respond with 204 No Content status code
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatCompletionRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqHead := r.Header.Get("Authorization")
	if !strings.HasPrefix(reqHead, "Bearer ") {
		http.Error(w, "Only Bearer authorization header is supported", http.StatusUnauthorized)
		return
	}
	//reqHead split
	splitRes := strings.Split(reqHead, "Bearer ")
	if len(splitRes) < 1 {
		http.Error(w, "Failed to split request header", http.StatusBadRequest)
		return
	}
	tokenStr := splitRes[1]

	jwtItem, found := c.Get(tokenStr)
	if found {
		jwt = jwtItem.(string)
	} else {
		jwt, err = auth.FetchJWT(splitRes[1])
		if err != nil {
			log.Error().Err(err).Msg("Error fetching JWT: ")
			http.Error(w, "Failed to fetch JWT", http.StatusInternalServerError)
			return
		}
		c.Set(tokenStr, jwt, 4*time.Minute)
	}

	orgIdItem, found := c.Get(tokenStr + ":orgId")
	if found {
		orgID = orgIdItem.(string)
	} else {
		orgID, err = user.FetchUserProfile(jwt)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching user profile: ")
			http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
			return
		}
		// 理论上access_token 对应的orgId是不会变的，所以可以一直缓存
		c.Set(tokenStr+":orgId", orgID, 24*time.Hour)
	}
	var response *http.Response
	if req.Stream {
		// while use stream
		response, err = stream.FetchStream(jwt, orgID, req.Messages, req.ModelType, req.MaxTokens) // Make sure to adjust the FetchStream function to return the response
		if err != nil {
			log.Error().Err(err).Msg("Error fetching stream: ")
			http.Error(w, "Failed to fetch stream", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
	} else {
		url := "https://api.groq.com/openai/v1/chat/completions"
		payload := map[string]interface{}{
			"model":       req.ModelType,
			"messages":    req.Messages,
			"temperature": 0.2,
			"max_tokens":  req.MaxTokens,
			"top_p":       0.8,
			"stream":      false,
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			log.Error().Err(err).Msg("Error creating request: ")
			return
		}
		req.Header.Set("Authorization", "Bearer "+jwt)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("groq-organization", orgID)
		client := &http.Client{}
		response, err = client.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Error sending request: ")
			return
		}
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Accel-Buffering", "no")
	utils.SetCorsHeaders(w)
	w.WriteHeader(http.StatusOK)
	buf := make([]byte, 4*1024)

	for {
		n, err := response.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
		}

		if err != nil {
			if err == io.EOF {
				//w.WriteHeader(http.StatusOK)
				err := response.Body.Close()
				if err != nil {
					log.Error().Err(err).Msg("Failed to close response")
					http.Error(w, "Failed to close response", http.StatusInternalServerError)
					return
				}
			}
			break
		}
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/Star-Studio-Develop/Groq2API", http.StatusFound)
}

func main() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(4*time.Minute, 1*time.Minute)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		chatCompletionsHandler(w, r, c)
	})

	log.Info().Msg("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
