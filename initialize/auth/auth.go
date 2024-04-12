package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TokenResponse struct {
	Data struct {
		SessionJWT string `json:"session_jwt"`
	} `json:"data"`
}

func FetchJWT(refreshToken string) (string, error) {
	url := "https://web.stytch.com/sdk/v1/sessions/authenticate"
	body := bytes.NewBuffer(nil)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+refreshToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("origin", "https://groq.com")
	req.Header.Set("referer", "https://groq.com/")
	req.Header.Add("x-sdk-client", "eyJldmVudF9pZCI6ImV2ZW50LWlkLWNkOWRiYzUzLTIzOTQtNDgxNy1hOTNlLTBkYzBiODlhNmQxNiIsImFwcF9zZXNzaW9uX2lkIjoiYXBwLXNlc3Npb24taWQtMGNmYjkxYzktZTMyNi00MDUxLWI1YzUtMDI5NjY2NjM3NDYzIiwicGVyc2lzdGVudF9pZCI6InBlcnNpc3RlbnQtaWQtYjM2MzhjNDEtMjEzNi00Yjc1LTkxOTAtNDczYTAyZWE2M2Y5IiwiY2xpZW50X3NlbnRfYXQiOiIyMDI0LTA0LTEwVDA5OjI0OjIxLjYyM1oiLCJ0aW1lem9uZSI6IkFzaWEvU2hhbmdoYWkiLCJzdHl0Y2hfdXNlcl9pZCI6InVzZXItbGl2ZS1hNjYxZGJjZS0yMWVmLTRiZGMtYjZlNC0wZDNmMGVlODhhM2YiLCJzdHl0Y2hfc2Vzc2lvbl9pZCI6InNlc3Npb24tbGl2ZS1kMzliZmZhMi03YjU2LTQ5MDctOWMwYS00N2U1MGM4N2Y5NmEiLCJhcHAiOnsiaWRlbnRpZmllciI6Imdyb3EuY29tIn0sInNkayI6eyJpZGVudGlmaWVyIjoiU3R5dGNoLmpzIEphdmFzY3JpcHQgU0RLIiwidmVyc2lvbiI6IjQuNS4zIn19")
	req.Header.Add("x-sdk-parent-host", "https://groq.com")

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
