package redis_test

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	publisher "github.com/nachoconques0/websockets_fun/internal/publisher/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisPublisher_PublishMessage_Success(t *testing.T) {
	rMock, mock := redismock.NewClientMock()

	mock.ExpectXAdd(&redis.XAddArgs{
		Stream: "messages",
		Values: map[string]interface{}{"body": "oliss"},
	}).SetVal("1-0")

	p := publisher.New(rMock, "messages")

	err := p.PublishMessage("oliss")
	assert.NoError(t, err)
}

func TestRedisPublisher_PublishMessage_Failure(t *testing.T) {
	rMock, mock := redismock.NewClientMock()

	mock.ExpectXAdd(&redis.XAddArgs{
		Stream: "messages",
		Values: map[string]interface{}{"body": "oliss"},
	}).SetErr(redis.Nil)

	p := publisher.New(rMock, "messages")

	err := p.PublishMessage("oliss")
	assert.Error(t, err)
}
