package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"houseparty.com/config"
	"houseparty.com/models"
)

type Tracks map[string]interface{}

func GetSongById(id string) (*models.Song, error) {
	token, err := config.GetSpotifyTokenObject()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/tracks/"+id, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New("spotify search request failed: " + string(body))
	}

	var responseBody Tracks
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	song, err := SimplifyTrack(responseBody)
	log.Println(song)
	return song, err

}

func SearchSongs(search string) (*Tracks, error) {
	token, err := config.GetSpotifyTokenObject()
	if err != nil {
		return nil, err
	}

	query := BuildQuery(search)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search?q="+query+"&type=track&limit=5", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New("spotify search request failed: " + string(body))
	}

	var responseBody Tracks
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil

}

func BuildQuery(search string) string {
	query := strings.Replace(search, " ", "+", -1)
	log.Println(query)
	query = strings.Trim(query, " ")

	return query
}

func SimplifyTracks(tracks *Tracks) ([]models.Song, error) {
	var songs []models.Song

	items := (*tracks)["tracks"].(map[string]interface{})["items"].([]interface{})

	for _, item := range items {
		track := item.(map[string]interface{})

		var artistsNames []string
		if artists, ok := track["artists"].([]interface{}); ok {
			for _, a := range artists {
				artist, ok := a.(map[string]interface{})
				if !ok {
					continue
				}
				artistsNames = append(artistsNames, artist["name"].(string))
			}
		}

		albumData, ok := track["album"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid album data")
		}
		albumName := albumData["name"].(string)

		var imageURL string
		var imageWidth, imageHeight int
		if images, ok := albumData["images"].([]interface{}); ok && len(images) > 0 {
			image := images[2].(map[string]interface{})
			imageURL = image["url"].(string)
			imageWidth = int(image["width"].(float64))
			imageHeight = int(image["height"].(float64))
		}

		id, _ := track["id"].(string)
		uri, _ := track["uri"].(string)
		name, _ := track["name"].(string)
		durationMs, _ := track["duration_ms"].(float64)
		explicit, _ := track["explicit"].(bool)
		externalURL, _ := track["external_urls"].(map[string]interface{})["spotify"].(string)

		songs = append(songs, models.Song{
			Id:      id,
			URI:     uri,
			Name:    name,
			Artists: artistsNames,
			Album:   albumName,
			Image: models.Image{
				URL:    imageURL,
				Width:  imageWidth,
				Height: imageHeight,
			},
			DurationMs:  int(durationMs),
			Explicit:    explicit,
			ExternalURL: externalURL,
		})
	}

	return songs, nil
}

func SimplifyTrack(track map[string]interface{}) (*models.Song, error) {
	var artistsNames []string
	if artists, ok := track["artists"].([]interface{}); ok {
		for _, a := range artists {
			artist, ok := a.(map[string]interface{})
			if !ok {
				continue
			}
			artistsNames = append(artistsNames, artist["name"].(string))
		}
	}

	albumData, ok := track["album"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid album data")
	}
	albumName := albumData["name"].(string)

	var imageURL string
	var imageWidth, imageHeight int
	if images, ok := albumData["images"].([]interface{}); ok && len(images) > 0 {
		image := images[2].(map[string]interface{})
		imageURL = image["url"].(string)
		imageWidth = int(image["width"].(float64))
		imageHeight = int(image["height"].(float64))
	}

	id, _ := track["id"].(string)
	uri, _ := track["uri"].(string)
	name, _ := track["name"].(string)
	durationMs, _ := track["duration_ms"].(float64)
	explicit, _ := track["explicit"].(bool)
	externalURL, _ := track["external_urls"].(map[string]interface{})["spotify"].(string)

	song := models.Song{
		Id:      id,
		URI:     uri,
		Name:    name,
		Artists: artistsNames,
		Album:   albumName,
		Image: models.Image{
			URL:    imageURL,
			Width:  imageWidth,
			Height: imageHeight,
		},
		DurationMs:  int(durationMs),
		Explicit:    explicit,
		ExternalURL: externalURL,
	}

	return &song, nil
}
