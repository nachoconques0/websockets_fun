package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"

	internalErrors "github.com/nachoconques0/websockets_fun/internal/errors"
)

var (
	ErrUpgradingConnection = internalErrors.NewInternalError("error upgrading HTTP connection")
	ErrPublishingMessage   = internalErrors.NewInternalError("error publishing message")
	ErrReadingMessage      = internalErrors.NewInternalError("error reading message")
	wsUpgrader             = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

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
		slog.Error(fmt.Sprintf("HandleIncomingConnection err: %s\n", err))
		responseError(w, r, ErrUpgradingConnection)
		return
	}
	c.service.AddClient(conn)
	defer c.service.RemoveClient(conn)

	slog.Info("connection received")
	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}

// responseError handles internals error http response.
func responseError(w http.ResponseWriter, r *http.Request, err error) {
	var internalErr *internalErrors.Error
	if errors.As(err, &internalErr) {
		w.WriteHeader(internalErr.HTTPStatus())
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(internalErr); err != nil {
			slog.Error(fmt.Sprintf("error encoding response: %s", err))
		}
		return
	}
	responseError(w, r, err)
}
