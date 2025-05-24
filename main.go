package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nachoconques0/websockets_fun/internal/broadcaster"
	subscriber "github.com/nachoconques0/websockets_fun/internal/broadcaster/controller"
	"github.com/nachoconques0/websockets_fun/internal/config"
	"github.com/nachoconques0/websockets_fun/internal/manager"
	redisPublisher "github.com/nachoconques0/websockets_fun/internal/publisher/redis"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()
	mux := http.NewServeMux()

	// ✅ Creamos el cliente Redis una sola vez
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	setupManager(cfg, mux, redisClient)
	setupBroadcaster(cfg, mux, redisClient)

	fmt.Println("Server running on", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func setupManager(cfg *config.Config, mux *http.ServeMux, redisClient *redis.Client) {
	rp := redisPublisher.New(redisClient, cfg.RedisStream)
	wsManager := manager.NewManager(rp)

	mux.HandleFunc("/ws/manager", wsManager.ServeWS)
	fmt.Println("WS Manager ready at /ws/manager")
}

func setupBroadcaster(cfg *config.Config, mux *http.ServeMux, redisClient *redis.Client) {
	b := broadcaster.New()
	ctrl := subscriber.New(redisClient, cfg.RedisStream, b)
	go ctrl.Start()

	mux.HandleFunc("/ws/broadcaster", b.ServeWS)
	fmt.Println("WS Broadcaster ready at /ws/broadcaster")
}
