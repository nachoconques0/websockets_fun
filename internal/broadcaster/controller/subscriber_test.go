package subscriber_test

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	subscriber "github.com/nachoconques0/websockets_fun/internal/broadcaster/controller"
	"github.com/nachoconques0/websockets_fun/internal/mocks"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

func TestController_Start_CallsBroadcast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBroadcaster := mocks.NewMockBroadcaster(ctrl)
	mockBroadcaster.EXPECT().
		Broadcast("hello-world").
		Times(1)

	db, mock := redismock.NewClientMock()

	mock.ExpectXRead(&redis.XReadArgs{
		Streams: []string{"messages", "0"},
		Block:   0,
	}).SetVal([]redis.XStream{
		{
			Stream: "messages",
			Messages: []redis.XMessage{
				{ID: "1-0", Values: map[string]interface{}{"body": "testinnnnnnnnnnnnnn"}},
			},
		},
	})

	// Crear el controller y lanzarlo
	c := subscriber.New(db, "messages", mockBroadcaster)
	go c.Start()

	time.Sleep(150 * time.Millisecond)
}
