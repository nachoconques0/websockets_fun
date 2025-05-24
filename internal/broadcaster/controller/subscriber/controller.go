package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Service contains service methods
type Service interface {
	AddClient(conn *websocket.Conn)
	RemoveClient(conn *websocket.Conn)
	Broadcast(msg string)
}

// Controllers holds the client and needed fields
type Controller struct {
	client  *redis.Client
	queue   string
	service Service
	ctx     context.Context
}

// New returns a new controller
func New(client *redis.Client, queue string, service Service) *Controller {
	return &Controller{
		client:  client,
		queue:   queue,
		service: service,
		ctx:     context.Background(),
	}
}

// Start starts a infinte loop in order to read from the stream
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
