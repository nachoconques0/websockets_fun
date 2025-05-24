package main

import (
	"fmt"
	"log"
	"net/http"

	broadcasterCtl "github.com/nachoconques0/websockets_fun/internal/broadcaster/controller/subscriber"
	broadcasterWSCtl "github.com/nachoconques0/websockets_fun/internal/broadcaster/controller/websocket"
	broadcasterSvc "github.com/nachoconques0/websockets_fun/internal/broadcaster/service/broadcaster"
	"github.com/nachoconques0/websockets_fun/internal/config"
	managerCtl "github.com/nachoconques0/websockets_fun/internal/manager/controller/manager"
	managerSvc "github.com/nachoconques0/websockets_fun/internal/manager/service/manager"
	redisPublisher "github.com/nachoconques0/websockets_fun/internal/publisher/redis"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()
	mux := http.NewServeMux()

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
	managerService := managerSvc.New(rp)
	managerController := managerCtl.New(managerService)

	mux.HandleFunc("/ws/manager", managerController.HandleIncomingConnection)
	fmt.Println("WS Manager ready at /ws/manager")
}

func setupBroadcaster(cfg *config.Config, mux *http.ServeMux, redisClient *redis.Client) {
	broadcasterService := broadcasterSvc.New()
	broadcasterController := broadcasterCtl.New(redisClient, cfg.RedisStream, broadcasterService)
	broadcasterWSController := broadcasterWSCtl.New(broadcasterService)

	go broadcasterController.Start()

	mux.HandleFunc("/ws/broadcaster", broadcasterWSController.HandleIncomingConnection)
	fmt.Println("WS Broadcaster ready at /ws/broadcaster")
}
