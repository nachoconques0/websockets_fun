package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Broadcaster interface {
	Broadcast(msg string)
}

type Controller struct {
	client      *redis.Client
	queue       string
	broadcaster Broadcaster
	ctx         context.Context
}

func New(client *redis.Client, queue string, broadcaster Broadcaster) *Controller {
	return &Controller{
		client:      client,
		queue:       queue,
		broadcaster: broadcaster,
		ctx:         context.Background(),
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
					c.broadcaster.Broadcast(text)
					fmt.Println("Broadcasted:", text)
				}
			}
		}
	}
}
