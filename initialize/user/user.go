package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserProfile struct {
	User struct {
		Orgs struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		} `json:"orgs"`
	} `json:"user"`
}

func FetchUserProfile(jwt string) (string, error) {
	url := "https://api.groq.com/platform/v1/user/profile"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	log.Printf("jwt: %v", jwt)
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var profile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return "", err
	}

	if len(profile.User.Orgs.Data) == 0 {
		return "", fmt.Errorf("no organizations found for user")
	}

	return profile.User.Orgs.Data[0].ID, nil
}
