package service

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Service struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]bool
}

func New() *Service {
	return &Service{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient registers a new WebSocket connection
func (s *Service) AddClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[conn] = true
}

// RemoveClient removes a WebSocket connection
func (s *Service) RemoveClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, conn)
	conn.Close()
}

// Broadcast sends a message to all active WebSocket clients
func (s *Service) Broadcast(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			// Clean up dead connection
			delete(s.clients, conn)
			conn.Close()
		}
	}
}
