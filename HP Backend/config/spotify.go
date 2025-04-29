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
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TimeIssued   int    `json:"time_issued"`
	UserID       int64  `json:"user_id"`
}

var redirectUrl = os.Getenv("FRONTEND_URL")

func (s *SpotifyTokenObject) SaveToken() error {

	deleteStmt, err := storage.DB.Prepare(storage.DeleteTokenQuery)
	if err != nil {
		return err
	}
	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(s.UserID)
	if err != nil {
		return err
	}

	stmt, err := storage.DB.Prepare(storage.SaveTokenQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.AccessToken, s.TokenType, s.Scope, s.ExpiresIn, s.RefreshToken, s.TimeIssued, s.UserID)

	if err != nil {
		return err
	}
	return nil
}

func (s *SpotifyTokenObject) UpdateToken() error {
	stmt, err := storage.DB.Prepare(storage.UpdateTokenQuery)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.AccessToken, s.TokenType, s.Scope, s.ExpiresIn, s.TimeIssued, s.UserID)
	if err != nil {
		return err
	}

	return nil
}

func GetTokenFromDB(hostId int64) (*SpotifyTokenObject, error) {
	var token SpotifyTokenObject
	row := storage.DB.QueryRow(storage.GetSpotifyToken, hostId)
	err := row.Scan(&token.AccessToken, &token.TokenType, &token.Scope, &token.ExpiresIn, &token.RefreshToken, &token.TimeIssued, &token.UserID)

	if err == sql.ErrNoRows {
		return nil, errors.New("token not found")
	} else if err != nil {
		return nil, err
	}

	return &token, nil
}

func checkIfTokenExpired(token *SpotifyTokenObject) bool {
	return token.TimeIssued+token.ExpiresIn < int((int64)(time.Now().Unix()))
}

func GetSpotifyTokenObject(hostId int64) (*SpotifyTokenObject, error) {
	token, err := GetTokenFromDB(hostId)
	if err != nil {
		log.Println("Error with getting from db")
		return nil, err
	}

	if checkIfTokenExpired(token) {
		tokenRefresh, err := RefreshToken(token.RefreshToken, token.UserID)
		if err != nil {
			return nil, err
		}
		token = tokenRefresh
	}

	return token, nil

}

func SetSpotifyToken(code string, userId int64) (*SpotifyTokenObject, error) {

	log.Println("Code Here: " + code)

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("missing spotify client id or client secret" + clientID + clientSecret)
	}

	data := url.Values{}
	data.Set("code", code)
	data.Add("redirect_uri", redirectUrl)
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
	tokenObject.UserID = userId

	err = tokenObject.SaveToken()
	if err != nil {
		return nil, err
	}

	return &tokenObject, nil
}

func GenerateSpotifyAuthRequest() (string, error) {

	scope := "streaming user-read-email user-read-private user-modify-playback-state user-read-playback-state"

	state := generateRandomString(16)
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")

	if clientID == "" {
		return "", errors.New("missing spotify client id or client secret" + clientID)
	}

	data := url.Values{}
	data.Add("response_type", "code")
	data.Add("client_id", clientID)
	data.Add("scope", scope)
	data.Add("redirect_uri", redirectUrl)
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

func RefreshToken(refreshToken string, userId int64) (*SpotifyTokenObject, error) {
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
	log.Println(resp.Body)
	tokenObject.TimeIssued = int((int64)(time.Now().Unix()))
	tokenObject.UserID = userId

	err = tokenObject.UpdateToken()
	if err != nil {
		return nil, err
	}
	return &tokenObject, nil
}
