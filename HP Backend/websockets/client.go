package websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"houseparty.com/models"
)

type ClientList map[*Client]bool
type PlayList []models.Song
type SkipRecord []int64

type Client struct {
	User       *models.User
	Connection *websocket.Conn
	RoomID     string
	Manager    *Manager
	Egress     chan Event
}

var (
	pongWait     = 60 * time.Second
	pingInterval = (pongWait * 9) / 10
)

func NewClient(connection *websocket.Conn, manager *Manager, RoomID string, userId int64) *Client {
	var user *models.User
	user.GetUserById(userId)

	return &Client{
		User: user,
		Connection: connection,
		RoomID:     RoomID,
		Manager:    manager,
		Egress:     make(chan Event),
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	err := c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println("failed to set read deadline: ", err.Error())
		return
	}

	c.Connection.SetReadLimit(512)
	c.Connection.SetPongHandler(c.PongHnadler)

	for {
		_, payLoad, err := c.Connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err.Error())
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payLoad, &request); err != nil {
			log.Println("failed to unmarshal message: ", err.Error())
			break
		}

		if err := c.Manager.routeEvent(request, c); err != nil {
			log.Println("failed to route event: ", err.Error())
			break
		}

	}
}

func (c *Client) WriteMessages() {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.Egress:

			if !ok {
				if err := c.Connection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Println("conncetion closed: ", err.Error())
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println("failed to marshal message: ", err.Error())
				return
			}

			if err := c.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("failed to send message: ", err.Error())
			}

		case <-ticker.C:
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("failed to send ping: ", err.Error())
				return
			}
		}
	}
}

func (c *Client) PongHnadler(pongMessage string) error {
	return c.Connection.SetReadDeadline(time.Now().Add(pongWait))
}
