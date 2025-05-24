# GIG Challenge
### Made by <3 Juan Calcagno AKA Nacho

---

### What I've done

This project implements two microservices that communicate through a Redis Stream:

- **Manager**: Handles incoming WebSocket connections, receives messages from clients, and publishes them to a Redis Stream.
- **Broadcaster**: Subscribes to the Redis Stream and broadcasts any incoming messages to all currently connected WebSocket clients.

### Postman Collection available :white_check_mark:
[Here should be the WS collection to test it](https://.postman.co/workspace/My-Workspace~8a3785f5-2a83-4cb8-8641-69b834d12c79/collection/682e2f60465421c338703047?action=share&creator=5221820&active-environment=5221820-991cd2df-2b04-4f2c-970c-1a909f5243f2)

### How to run it :scream_cat:
1. Have docker in your machine
2. `git clone` this repo
3. Once you are inside the repo
4. Run `docker compose up -d` this will initiate a container with a running redis env
5. Run `make mod` so you download needed pkgs
6. Run `make run` and if all good. Project should be running ready to get some WS connections
7. Connect with one client to the manager `localhost:3000/ws/manager`
8. Connect with another client to the broadcaster `localhost:3000/ws/broadcaster`
9. From the client tab that is connected to the manager send a message. This should appear in the client tab connected to the broadcaster

### You don't want to run it? :smiling_imp:
1. Have docker in your machine
2. `git clone` this repo
3. Once you are inside the repo
4. Run `make mock` 
5. Run `make test` and this will trigger a docker compose file that will spin up redis container and then run all the needed tests. By the time of writing this test are passing lol. 

### If you wanna clean the Redis stream. Run this command
`docker exec -it websockets_fun-redis-1 redis-cli XTRIM local-q MAXLEN 0`

### WebSocket endpoints
	ws://localhost:8080/ws/manager → For sending messages.

	ws://localhost:8080/ws/broadcast → For receiving broadcasts.

### Structure :palm_tree:
```
📦websockets_fun
 ┣ 📂internal
 ┃ ┣ 📂broadcaster
 ┃ ┃ ┣ 📂controller
 ┃ ┃ ┃ ┣ 📂subscriber
 ┃ ┃ ┃ ┃ ┣ 📜controller.go
 ┃ ┃ ┃ ┃ ┗ 📜controller_test.go
 ┃ ┃ ┃ ┗ 📂websocket
 ┃ ┃ ┃ ┃ ┣ 📜controller.go
 ┃ ┃ ┃ ┃ ┗ 📜controller_test.go
 ┃ ┃ ┗ 📂service
 ┃ ┃ ┃ ┗ 📂broadcaster
 ┃ ┃ ┃ ┃ ┣ 📜service.go
 ┃ ┃ ┃ ┃ ┗ 📜service_test.go
 ┃ ┣ 📂config
 ┃ ┃ ┗ 📜config.go
 ┃ ┣ 📂errors
 ┃ ┃ ┗ 📜errors.go
 ┃ ┣ 📂manager
 ┃ ┃ ┣ 📂controller
 ┃ ┃ ┃ ┗ 📂manager
 ┃ ┃ ┃ ┃ ┣ 📜controller.go
 ┃ ┃ ┃ ┃ ┗ 📜controller_test.go
 ┃ ┃ ┗ 📂service
 ┃ ┃ ┃ ┗ 📂manager
 ┃ ┃ ┃ ┃ ┣ 📜service.go
 ┃ ┃ ┃ ┃ ┗ 📜service_test.go
 ┃ ┣ 📂mocks
 ┃ ┃ ┣ 📜mock_broadcaster_controller.go
 ┃ ┃ ┣ 📜mock_manager_controller.go
 ┃ ┃ ┗ 📜mock_manager_service.go
 ┃ ┗ 📂publisher
 ┃ ┃ ┗ 📂redis
 ┃ ┃ ┃ ┣ 📜publisher.go
 ┃ ┃ ┃ ┗ 📜publisher_test.go
 ┣ 📜,gitignore
 ┣ 📜DISCLAIMERS.md
 ┣ 📜Makefile
 ┣ 📜README.md
 ┣ 📜docker-compose.yml
 ┣ 📜generate-mocks.sh
 ┣ 📜go.mod
 ┣ 📜go.sum
 ┗ 📜main.go
```
