package broadcaster

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Broadcaster struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func New() *Broadcaster {
	return &Broadcaster{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (b *Broadcaster) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket broadcaster upgrade failed:", err)
		return
	}
	b.AddClient(conn)
	defer b.RemoveClient(conn)

	// Keep the connection open until the client disconnects
	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}

func (b *Broadcaster) AddClient(conn *websocket.Conn) {
	fmt.Println("adding client?")
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[conn] = true
}

func (b *Broadcaster) RemoveClient(conn *websocket.Conn) {
	fmt.Println("removing client?")
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, conn)
	conn.Close()
}

func (b *Broadcaster) Broadcast(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for conn := range b.clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			b.RemoveClient(conn)
		}
	}
}
