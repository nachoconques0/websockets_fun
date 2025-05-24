package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Service holds needed methods for the controller to be used
type Service interface {
	PublishMessage(msg string) error
}

// Controller holds service
type Controller struct {
	Service Service
}

// New returns a new controller
func New(service Service) *Controller {
	return &Controller{
		Service: service,
	}
}

// HandleIncomingConnection handles incoming connection, upgrated it to use WS and then
// calls service to publish to queue
func (c *Controller) HandleIncomingConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("there was an error upgrading the HTTP connection")
		return
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("there was an error reading message:", err)
			break
		}
		err = c.Service.PublishMessage(string(msg))
		if err != nil {
			fmt.Println("there was an error publishing message")
		}
	}
}
