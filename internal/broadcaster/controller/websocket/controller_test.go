package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/gorilla/websocket"
	controller "github.com/nachoconques0/websockets_fun/internal/broadcaster/controller/websocket"
	"github.com/nachoconques0/websockets_fun/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleIncomingConnection_AddAndRemoveClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockBroadcasterService(ctrl)
	defer ctrl.Finish()

	mockService.EXPECT().AddClient(gomock.Any()).Times(1)
	mockService.EXPECT().RemoveClient(gomock.Any()).Times(1)

	wsController := controller.New(mockService)

	server := httptest.NewServer(http.HandlerFunc(wsController.HandleIncomingConnection))
	defer server.Close()

	url := "ws" + server.URL[4:]

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	conn.Close()

	time.Sleep(100 * time.Millisecond)
}
