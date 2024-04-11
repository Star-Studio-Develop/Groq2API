package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"Groq2API/initialize/auth"
	"Groq2API/initialize/model"
	"Groq2API/initialize/stream"
	"Groq2API/initialize/user"
)

type ChatCompletionRequest struct {
	//RefreshToken string `json:"refresh_token"`
	Messages  []model.Message `json:"messages"`
	ModelType string          `json:"model"`
}

func chatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	jwt, err := auth.FetchJWT(splitRes[1])
	if err != nil {
		log.Printf("Error fetching JWT: %v", err)
		http.Error(w, "Failed to fetch JWT", http.StatusInternalServerError)
		return
	}

	orgID, err := user.FetchUserProfile(jwt)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	response, err := stream.FetchStream(jwt, orgID, req.Messages, req.ModelType) // Make sure to adjust the FetchStream function to return the response instead of printing it.
	if err != nil {
		log.Printf("Error fetching stream: %v", err)
		http.Error(w, "Failed to fetch stream", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Accel-Buffering", "no")
	//w.WriteHeader(http.StatusOK)
	buf := make([]byte, 4*1024)

	for {
		n, err := response.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
		}

		if err != nil {
			if err == io.EOF {
				w.WriteHeader(http.StatusOK)
				err := response.Body.Close()
				if err != nil {
					log.Printf("Error close response: %v", err)
					http.Error(w, "Failed to close response", http.StatusInternalServerError)
					return
				}
			}
			break
		}
	}

}
func main() {
	http.HandleFunc("/v1/chat/completions", chatCompletionsHandler)

	fmt.Println("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
