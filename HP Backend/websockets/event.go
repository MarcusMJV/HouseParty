package websockets

import (
	"encoding/json"
	"log"
	"slices"
	"time"

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
	EventSkipRequest  = "skip-request"
	EventSongSkipped  = "song-skipped"
	FinalSongEnded    = "final-song-ended"
	UserLeft          = "user-left"
)

// Define Event Struct and Event Handler
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type EventHandler func(event Event, c *Client) error

// Define Event Payloads
type JoinedRoomEvent struct {
	UserCount    int           `json:"user_count"`
	PlayList     []models.Song `json:"playlist"`
	CurrentSong  *models.Song  `json:"current_song"`
	ApiToken     string        `json:"api_token"`
	SongPosition int64         `json:"song_position"`
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
	room := c.Manager.Rooms[c.RoomID]

	for client := range room.Clients {
		if client == c {
			continue
		}
		client.Egress <- event
	}

	apiToken, err := config.GetSpotifyTokenObject(room.HostID)
	if err != nil {
		return err
	}
	songPosition := time.Since(room.CurrentSongStartedAt)

	joinedEvent := JoinedRoomEvent{
		UserCount:    len(room.Clients),
		PlayList:     room.PlayList,
		CurrentSong:  room.CurrentSong,
		ApiToken:     apiToken.AccessToken,
		SongPosition: songPosition.Milliseconds(),
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
	room := c.Manager.Rooms[c.RoomID]

	if err != nil {
		log.Println(err)
		return err
	}
	tracks, err := services.SearchSongs(searchEvent.Search, room.HostID)
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

	song, err := services.GetSongById(addSongEvent.SongId, room.HostID)
	if err != nil {
		return err
	}

	if room.CurrentSong == nil {
		err := room.PrepareSongToPlay(song)
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

func SkipSongRequest(event Event, c *Client) error {
	room := c.Manager.Rooms[c.RoomID]
	if slices.Contains(room.UserSkipRecord[:], c.User.Id) {
		return nil
	}

	response := event

	threshold := len(room.Clients) / 2
	room.UserSkipRecord = append(room.UserSkipRecord, c.User.Id)

	if len(room.UserSkipRecord) > threshold {
		room.SkipChan <- true
		room.UserSkipRecord = SkipRecord{}

		return nil
	}
	room.UserSkipRecord = append(room.UserSkipRecord, c.User.Id)

	room.SendEventToAllClients(response)
	return nil
}

func HandleUserLeaving(event Event, c *Client) error {
	room := c.Manager.Rooms[c.RoomID]
	room.UserLeftRoom()
	return nil
}
