package websockets

import (
	"encoding/json"
	"log"

	"houseparty.com/config"
	"houseparty.com/models"
	"houseparty.com/services"
)

// Define Event Types
const (
	EventJoinRoom     = "joined-room"
	EventSearchSongs  = "search-songs"
	EventAddSong      = "add-song"
	EventPlaySong     = "play-song"
	AddedSongPlaylist = "added-song-playlist"
	SetAndPlaySong    = "set-and-play-song"
)

// Define Event Struct and Event Handler
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type EventHandler func(event Event, c *Client) error

// Define Event Payloads
type JoinedRoomEvent struct {
	UserCount   int           `json:"user_count"`
	PlayList    []models.Song `json:"playlist"`
	CurrentSong *models.Song  `json:"current_song"`
	ApiToken    string        `json:"api_token"`
}
type SongChangeEvent struct {
	PlayList    []models.Song `json:"playlist"`
	CurrentSong *models.Song  `json:"current_song"`
	ApiToken    string        `json:"api_token"`
}

type SearchSongsEvent struct {
	Search string `json:"search"`
}
type SearchResultsEvent struct {
	Songs []models.Song `json:"songs"`
}

type AddSongEvent struct {
	From   string `json:"from"`
	SongId string `json:"song_id"`
}
type AddSongResultEvent struct {
	From     string       `json:"from"`
	Song     *models.Song `json:"song"`
	ApiToken string       `json:"api_token"`
}

type AddedSongToPlaylist struct {
	From string       `json:"from"`
	Song *models.Song `json:"song"`
}

type SetAndPlayCurrentSong struct {
	ApiToken string       `json:"api_token"`
	Song     *models.Song `json:"song"`
}

// Define Event Handlers
func JoinRoom(event Event, c *Client) error {
	for client := range c.Manager.Rooms[c.RoomID].Clients {
		if client == c {
			continue
		}
		client.Egress <- event
	}

	apiToken, err := config.GetSpotifyTokenObject()
	if err != nil {
		return err
	}

	joinedEvent := JoinedRoomEvent{
		UserCount:   len(c.Manager.Rooms[c.RoomID].Clients),
		PlayList:    c.Manager.Rooms[c.RoomID].PlayList,
		CurrentSong: c.Manager.Rooms[c.RoomID].CurrentSong,
		ApiToken:    apiToken.AccessToken,
	}

	joinedPayload, err := json.Marshal(joinedEvent)
	if err != nil {
		return err
	}

	event = Event{
		Type:    "room-information",
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
		Type:    EventSearchSongs,
		Payload: responsePayload,
	}

	c.Egress <- response
	return nil
}

func AddSong(event Event, c *Client) error {
	var addSongEvent AddSongEvent
	var response Event

	room := c.Manager.Rooms[c.RoomID]

	err := json.Unmarshal(event.Payload, &addSongEvent)
	if err != nil {
		return err
	}

	song, err := services.GetSongById(addSongEvent.SongId)
	if err != nil {
		return err
	}

	if room.CurrentSong == nil {
		err := room.SetCurrentSong(song)
		if err != nil {
			return err
		}

	} else {
		payload, err := room.AddSongToPlaylist(song, addSongEvent.From)
		if err != nil {
			return err
		}
		response.Type = AddedSongPlaylist
		response.Payload = payload
		room.SendEventToAllClients(response)
	}

	return nil
}
