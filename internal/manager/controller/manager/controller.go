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
		slog.Error(fmt.Sprintf("HandleIncomingConnection err: %s\n", err))
		responseError(w, r, ErrUpgradingConnection)
		return
	}
	defer conn.Close()

	slog.Info("connection received")
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error(fmt.Sprintf("HandleIncomingConnection reading message err: %s\n", err))
			responseError(w, r, ErrReadingMessage)
			break
		}
		err = c.Service.PublishMessage(string(msg))
		if err != nil {
			slog.Error(fmt.Sprintf("HandleIncomingConnection publishing msg err: %s\n", err))
			responseError(w, r, ErrPublishingMessage)
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
