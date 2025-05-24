package broadcaster_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nachoconques0/websockets_fun/internal/broadcaster"
	"github.com/stretchr/testify/assert"
)

func TestBroadcaster_Broadcast_Success(t *testing.T) {
	b := broadcaster.New()
	server := httptest.NewServer(http.HandlerFunc(b.ServeWS))
	defer server.Close()

	url := "ws" + server.URL[4:] + "/ws/broadcaster"

	conn1, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)

	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)

	defer conn2.Close()

	time.Sleep(100 * time.Millisecond) // wait for connections to register

	testMsg := "olis"
	b.Broadcast(testMsg)

	for _, conn := range []*websocket.Conn{conn1, conn2} {
		_, msg, err := conn.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, string(msg), testMsg)
	}
}
