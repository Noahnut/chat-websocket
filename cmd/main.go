package main

import (
	"chat-websocket/api"
	"context"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	websocketServer := api.NewWebSocketServer(ctx, cancel, 8080)

	websocketServer.Run()
}
