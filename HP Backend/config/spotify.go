package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type AuthResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
	TimeIssued 	int    `json:"time_issued"`
}

var SpotifyToken *AuthResponse

func GetSpotifyToken() (*AuthResponse, error) {
	if CheckToken() {
		return SpotifyToken, nil
	}
	
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, errors.New("missing spotify client id or client secret" + clientID + clientSecret)
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
        return nil, err
    }

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
        return nil, err
    }

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("spotify token request failed: " + string(body))
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, err
	}
	authResponse.TimeIssued = int((int64)(time.Now().Unix()))

	SpotifyToken = &authResponse
	return &authResponse, nil
}

func CheckToken() bool {
	if SpotifyToken == nil || SpotifyToken.TimeIssued + SpotifyToken.ExpiresIn < int((int64)(time.Now().Unix())) {
		return false
	}
	return true
}