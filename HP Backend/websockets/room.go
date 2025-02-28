package websockets

import (
	"encoding/json"
	"time"
	"houseparty.com/config"
	"houseparty.com/models"
)
 
type RoomDataList map[string]*RoomData

type RoomData struct {
	*models.Room
	Clients ClientList
	PlayList PlayList
	CurrentSong *models.Song
	CurrentSongStartedAt time.Time
}

func NewRoomData(room *models.Room) *RoomData {
	var roomPlaylist []models.Song
	return &RoomData{
		Room: room,
		Clients: make(ClientList),
		PlayList: roomPlaylist,
		CurrentSong: nil,
	}
}

func (r *RoomData) SendEventToAllClients(event Event) {
	for client := range r.Clients {
		client.Egress <- event
	}
}

func (r *RoomData) AddSongToPlaylist(song *models.Song, name string) ([]byte, error) {
	r.PlayList = append(r.PlayList, *song)
	response := AddedSongToPlaylist{
		From: name,
		Song: song,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *RoomData) SetCurrentSong(song *models.Song) error {
	r.CurrentSong = song
	tokenObject , err:= config.GetSpotifyTokenObject(r.HostID)

	if err != nil {
		return err
	}

	response := SetAndPlayCurrentSong{
		ApiToken: tokenObject.AccessToken,
		Song: song,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	r.PlaySong(song, payload)
	return nil
}

func (r *RoomData) PlaySong(song *models.Song, payload []byte)  {
	r.CurrentSong = song
	timer := time.NewTimer(time.Duration(song.DurationMs) * time.Millisecond)

	event := Event{
		Type: SetAndPlaySong,
		Payload: payload,
	}

	r.SendEventToAllClients(event)
	r.CurrentSongStartedAt = time.Now()

	go func() {
		<-timer.C
		r.HandleSongFinished()
	}()
}

func (r *RoomData) HandleSongFinished() {
	if len(r.PlayList) > 0 {
		nextSong := r.PlayList[0]
		r.PlayList = r.PlayList[1:]

		r.SetCurrentSong(&nextSong)
	}
}
