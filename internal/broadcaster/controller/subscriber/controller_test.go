package controller_test

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"

	controller "github.com/nachoconques0/websockets_fun/internal/broadcaster/controller/subscriber"
	"github.com/nachoconques0/websockets_fun/internal/mocks"
)

func TestSubscriber_Listen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockBroadcasterService(ctrl)
	mockService.EXPECT().Broadcast("test-message").Times(1)

	rdMock, mock := redismock.NewClientMock()

	mock.ExpectXRead(&redis.XReadArgs{
		Streams: []string{"test-stream", "0"},
		Block:   0,
	}).SetVal([]redis.XStream{
		{
			Stream: "test-stream",
			Messages: []redis.XMessage{
				{ID: "1-0", Values: map[string]interface{}{"body": "test-message"}},
			},
		},
	})

	sub := controller.New(rdMock, "test-stream", mockService)
	go sub.Start()

	time.Sleep(150 * time.Millisecond)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}
