package user

import (
	"encoding/json"
	"fmt"
	"github.com/Star-Studio-Develop/Groq2API/initialize/model"
	"github.com/rs/zerolog/log"
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

func FetchUserProfile(refreshToken string, jwt string) (string, error) {

	if userCache, ok := model.GetUserCache(refreshToken); ok && userCache.OrgID != "" {
		log.Info().Msg("Using cached orgID")
		return userCache.OrgID, nil
	}

	log.Info().Msg("Fetching new orgID")
	url := "https://api.groq.com/platform/v1/user/profile"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	log.Info().Str("jwt", jwt).Msg("jwt token")
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

	// 更新缓存
	model.SetOrgID(refreshToken, profile.User.Orgs.Data[0].ID)

	return profile.User.Orgs.Data[0].ID, nil
}
