package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// Publisher contains Publisher fields
type Publisher struct {
	client *redis.Client
	ctx    context.Context
	stream string
}

// New returns a new Publisher client
func New(client *redis.Client, stream string) *Publisher {
	return &Publisher{
		client: client,
		ctx:    context.Background(),
		stream: stream,
	}
}

// PublishMessage publish a messsage to the respective stream
func (p *Publisher) PublishMessage(msg string) error {
	err := p.client.XAdd(p.ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]interface{}{"body": msg},
	}).Err()
	if err != nil {
		slog.Error(fmt.Sprintf("PublishMessage to redis err: %s\n", err))
		return err
	}
	fmt.Println("Published to Redis:", msg)
	return nil
}
