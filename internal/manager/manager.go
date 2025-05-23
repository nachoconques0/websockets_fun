package manager

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

type Publisher interface {
	PublishMessage(msg string) error
}

type Manager struct {
	publisher Publisher
}

func NewManager(p Publisher) *Manager {
	return &Manager{
		publisher: p,
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("We just got a connection")
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("there was an error upgrading the HTTP connection")
		return
	}
	defer conn.Close()

	m.publishIncomingMsg(conn)
}

func (m *Manager) publishIncomingMsg(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("there was an error reading message:", err)
			break
		}
		err = m.publisher.PublishMessage(string(msg))
		if err != nil {
			fmt.Println("there was an error publishing message")
		}
	}
}
