package websockets

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"houseparty.com/models"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}
)

type Manager struct {
	Rooms RoomDataList
	sync.RWMutex
	Handlers map[string]EventHandler
}

func (m *Manager) CountClients(roomID string) int {
	m.RLock()
	defer m.RUnlock()
	return len(m.Rooms[roomID].Clients)
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.Handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
	} else {
		return errors.New("no handler for event type")
	}
	return nil
}

func NewManager() *Manager {
	m := &Manager{
		Rooms:    make(RoomDataList),
		Handlers: make(map[string]EventHandler),
	}

	m.SetupEventHandlers()
	return m
}

func (m *Manager) SetupEventHandlers() {
	m.Handlers[EventJoinRoom] = JoinRoom
	m.Handlers[EventSearchSongs] = SearchSongs
	m.Handlers[EventAddSong] = AddSong
	m.Handlers[EventSkipRequest] = SkipSongRequest
	m.Handlers[UserLeft] = HandleUserLeaving
}

func (m *Manager) AddClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if m.Rooms[client.RoomID] == nil {
		room := &models.Room{}
		room.GetRoomById(client.RoomID)
		m.Rooms[client.RoomID] = NewRoomData(room)
	}

	m.Rooms[client.RoomID].Clients[client] = true
}

func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	room := m.Rooms[client.RoomID]

	if _, ok := room.Clients[client]; ok {
		client.Connection.Close()
		delete(m.Rooms[client.RoomID].Clients, client)
	}

}

func (m *Manager) ServeWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("New Connection")

		roomId := c.Param("id")

		conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		log.Println("USER ID IS HERE: ", c.GetInt64("userId"))

		client := NewClient(conn, m, roomId, c.GetInt64("userId"))
		m.AddClient(client)

		go client.ReadMessages()
		go client.WriteMessages()
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:5173":
		return true
	default:
		return false
	}
}
