package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	client *redis.Client
	ctx    context.Context
	stream string
}

func New(addr, stream string) *RedisPublisher {
	return &RedisPublisher{
		client: redis.NewClient(&redis.Options{Addr: addr}),
		ctx:    context.Background(),
		stream: stream,
	}
}

func (p *RedisPublisher) PublishMessage(msg string) error {
	err := p.client.XAdd(p.ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]interface{}{"body": msg},
	}).Err()
	if err != nil {
		fmt.Println("Error publishing to Redis:", err)
		return err
	}
	fmt.Println("Published to Redis:", msg)
	return nil
}
