package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	controller "github.com/nachoconques0/websockets_fun/internal/manager/controller/manager"
	"github.com/nachoconques0/websockets_fun/internal/mocks"
)

func TestHandleIncomingConnection_PublishesMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		PublishMessage("hello").
		Return(nil).
		Times(1)

	wsController := controller.New(mockService)

	server := httptest.NewServer(http.HandlerFunc(wsController.HandleIncomingConnection))
	defer server.Close()

	url := "ws" + server.URL[4:]

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
}
