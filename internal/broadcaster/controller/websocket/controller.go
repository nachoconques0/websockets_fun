package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Service defines the interface expected from the broadcaster service
type Service interface {
	AddClient(conn *websocket.Conn)
	RemoveClient(conn *websocket.Conn)
	Broadcast(msg string)
}

// Controller handles WebSocket connections
type Controller struct {
	service Service
}

// New creates a new WebSocket controller
func New(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

// HandleIncomingConnection upgrades the HTTP request to a WebSocket
// and manages the client lifecycle
func (c *Controller) HandleIncomingConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	c.service.AddClient(conn)
	defer c.service.RemoveClient(conn)

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}
