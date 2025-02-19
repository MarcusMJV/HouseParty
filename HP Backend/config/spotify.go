package config

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"houseparty.com/storage"
)

type SpotifyTokenObject struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	RefereshToken string `json:"refresh_token"`
	Scope         string `json:"scope"`
	TimeIssued    int    `json:"time_issued"`
}

var SpotifyToken *SpotifyTokenObject

func (s *SpotifyTokenObject) SaveToken() error {
	stmt, err := storage.DB.Prepare(storage.SaveTokenQuery);
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.AccessToken, s.TokenType, s.ExpiresIn, s.RefereshToken, s.Scope, s.TimeIssued)
	if err != nil {
		return err
	}
	return nil
}

func GetTokenFromDB() (*SpotifyTokenObject, error) {
	var token SpotifyTokenObject
	row := storage.DB.QueryRow(storage.GetSpotifyToken)
	err := row.Scan(&token.AccessToken, &token.TokenType, &token.ExpiresIn, &token.RefereshToken, &token.Scope, &token.TimeIssued)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &token, nil
}

func checkIfTokenExpired() bool {
	return SpotifyToken.TimeIssued+SpotifyToken.ExpiresIn < int((int64)(time.Now().Unix()))
}

func GetSpotifyTokenObject() (*SpotifyTokenObject, error) {
	if SpotifyToken == nil {
		token, err := GetTokenFromDB()
		if err != nil {
			return nil, err
		}
		SpotifyToken = token
	}

	if(checkIfTokenExpired()){
		token, err := RefereshToken(SpotifyToken.RefereshToken)
		if err != nil {
			return nil, err
		}
		SpotifyToken = token
	}

	return SpotifyToken, nil
	
}

func SetSpotifyToken(code string) (*SpotifyTokenObject, error) {

	log.Println("Code Here: " + code)

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("missing spotify client id or client secret" + clientID + clientSecret)
	}

	data := url.Values{}
	data.Set("code", code)
	data.Add("redirect_uri", "http://192.168.3.3:8080/spotify/token/callback")
	data.Add("grant_type", "authorization_code")

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	credentials := clientID + ":" + clientSecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Add("Authorization", "Basic "+encodedCredentials)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	var tokenObject SpotifyTokenObject
	if err := json.NewDecoder(resp.Body).Decode(&tokenObject); err != nil {
		return nil, err
	}
	tokenObject.TimeIssued = int((int64)(time.Now().Unix()))

	err = tokenObject.SaveToken()
	if err != nil {
		return nil, err
	}
	SpotifyToken = &tokenObject
	
	return &tokenObject, nil
}

func GenerateSpotifyAuthRequest() (string, error) {
	redirectURL := "http://192.168.3.3:8080/spotify/token/callback"
	scope := "streaming user-read-email user-read-private"
	state := generateRandomString(16)
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")

	if clientID == "" {
		return "", errors.New("missing spotify client id or client secret" + clientID)
	}

	data := url.Values{}
	data.Add("response_type", "code")
	data.Add("client_id", clientID)
	data.Add("scope", scope)
	data.Add("redirect_uri", redirectURL)
	data.Add("state", state)

	authUrl := "https://accounts.spotify.com/authorize/?" + data.Encode()

	return authUrl, nil
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("Error generating random string: ", err)
	}
	return hex.EncodeToString(b)[:n]
}

func RefereshToken(refreshToken string) (*SpotifyTokenObject, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("missing spotify client id or client secret" + clientID + clientSecret)
	}

	data := url.Values{}
	data.Add("refresh_token", refreshToken)
	data.Add("grant_type", "refresh_token")

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	credentials := clientID + ":" + clientSecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Add("Authorization", "Basic "+encodedCredentials)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	var tokenObject SpotifyTokenObject
	if err := json.NewDecoder(resp.Body).Decode(&tokenObject); err != nil {
		return nil, err
	}
	tokenObject.TimeIssued = int((int64)(time.Now().Unix()))

	err = tokenObject.SaveToken()
	if err != nil {
		return nil, err
	}
	return &tokenObject, nil
}

