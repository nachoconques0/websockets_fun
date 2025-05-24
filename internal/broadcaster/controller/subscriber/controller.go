package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	Broadcast(msg string)
}

type Controller struct {
	client  *redis.Client
	queue   string
	service Service
	ctx     context.Context
}

func New(client *redis.Client, queue string, service Service) *Controller {
	return &Controller{
		client:  client,
		queue:   queue,
		service: service,
		ctx:     context.Background(),
	}
}

func (c *Controller) Start() {
	lastID := "0"
	for {
		streams, err := c.client.XRead(c.ctx, &redis.XReadArgs{
			Streams: []string{c.queue, lastID},
			Block:   0,
		}).Result()

		if err != nil {
			fmt.Println("XREAD error:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				lastID = msg.ID
				if val, ok := msg.Values["body"]; ok {
					text := fmt.Sprintf("%v", val)
					c.service.Broadcast(text)
					fmt.Println("Broadcasted:", text)
				}
			}
		}
	}
}
