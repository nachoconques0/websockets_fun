package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nachoconques0/websockets_fun/internal/broadcaster"
	"github.com/nachoconques0/websockets_fun/internal/broadcaster/controller/subscriber"
	"github.com/nachoconques0/websockets_fun/internal/config"
	"github.com/nachoconques0/websockets_fun/internal/manager"
	"github.com/nachoconques0/websockets_fun/internal/publisher/redis"
)

func main() {
	cfg := config.Load()
	mux := http.NewServeMux()

	setupManager(cfg, mux)
	setupBroadcaster(cfg, mux)

	fmt.Println("Server running on", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func setupManager(cfg *config.Config, mux *http.ServeMux) {
	redisPublisher := redis.New(cfg.RedisAddr, cfg.RedisStream)
	wsManager := manager.NewManager(redisPublisher)

	mux.HandleFunc("/ws/manager", wsManager.ServeWS)
	fmt.Println("WS Manager ready at /ws/manager")
}

func setupBroadcaster(cfg *config.Config, mux *http.ServeMux) {
	b := broadcaster.New()
	ctrl := subscriber.New(cfg.RedisAddr, cfg.RedisStream, b)
	go ctrl.Start()

	mux.HandleFunc("/ws/broadcaster", b.ServeWS)
	fmt.Println("WS Broadcaster ready at /ws/broadcaster")
}
