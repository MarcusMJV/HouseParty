package websockets

import (
	"encoding/json"
	"log"

	"houseparty.com/config"
	"houseparty.com/models"
	"houseparty.com/services"
)

//Define Event Types
const(
	EventJoinRoom = "joined-room"
	EventSearchSongs = "search-songs"
	EventAddSong = "add-song"
)

//Define Event Struct and Event Handler
type Event struct {
	Type    string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type EventHandler func(event Event, c *Client) error

//Define Event Payloads
type JoinedRoomEvent struct {
	UserCount int `json:"user_count"`
	PlayList []models.Song `json:"playlist"`
	CurrentSong *models.Song `json:"current_song"`
}

type SearchSongsEvent struct {
	Search string `json:"search"`
}
type SearchResultsEvent struct {
	Songs []models.Song `json:"songs"`
}

type AddSongEvent struct {
	From string `json:"from"`
	SongId string `json:"song_id"`
}
type AddSongResultEvent struct {
	From string `json:"from"`
	Song *models.Song `json:"song"`
	ApiToken string `json:"api_token"`
}

//Define Event Handlers
func JoinRoom(event Event, c *Client) error {
	for client := range c.Manager.Rooms[c.RoomID].Clients {
		if client == c {
			continue
		}
		client.Egress <- event
	}

	joinedEvent := JoinedRoomEvent{
		UserCount: len(c.Manager.Rooms[c.RoomID].Clients),
		PlayList: c.Manager.Rooms[c.RoomID].PlayList,
		CurrentSong: &c.Manager.Rooms[c.RoomID].CurrentSong,
	}

	joinedPayload, err := json.Marshal(joinedEvent)
	if err != nil {
		return err
	}

	event = Event{
		Type: "room-information",
		Payload: joinedPayload,
	}
	c.Egress <- event
	return nil
}

func SearchSongs(event Event, c *Client) error {
	var searchEvent SearchSongsEvent
	err := json.Unmarshal(event.Payload, &searchEvent)
	if err != nil {
		log.Println(err)
		return err
	}
	tracks, err := services.SearchSongs(searchEvent.Search)
	if err != nil {
		log.Println(err)
		return err
	}
	songs, err := services.SimplifyTracks(tracks)
	if err != nil {
		log.Println(err)
		return err
	}

	responsePayload, err := json.Marshal(SearchResultsEvent{Songs: songs})
	if err != nil {
		return err
	}

	response := Event{
		Type: EventSearchSongs,
		Payload: responsePayload,
	}

	c.Egress <- response
	return nil
}

func AddSong(event Event, c *Client) error {
	var addSongEvent AddSongEvent
	err := json.Unmarshal(event.Payload, &addSongEvent)
	if err != nil {
		return err
	}

	song, err := services.GetSongById(addSongEvent.SongId)
	if err != nil {
		return err
	}

	log.Println(song)

	c.Manager.AddSongToPlaylist(song, c.RoomID)

	apiToken, err := config.GetSpotifyToken()
	if err != nil {
		return err
	}

	responsePayload, err := json.Marshal(AddSongResultEvent{From: addSongEvent.From, Song: song, ApiToken: apiToken.AccessToken})
	if err != nil {
		return err
	}

	event = Event{
		Type: EventAddSong,
		Payload: responsePayload,
	}
	for client := range c.Manager.Rooms[c.RoomID].Clients {
		client.Egress <- event
	}
	return nil
}
