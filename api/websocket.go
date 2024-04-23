package api

import (
	"chat-websocket/internal"
	"chat-websocket/queue"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type WebSocketServer struct {
	port          int
	websocketConn *internal.WebSocketConns
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewWebSocketServer(ctx context.Context, cancel context.CancelFunc, port int) *WebSocketServer {

	queue := queue.NewNATS("nats://192.168.0.109:30176")

	if err := queue.Connect(); err != nil {
		panic(err)
	}

	return &WebSocketServer{
		port:          port,
		websocketConn: internal.NewWebSocketConns(1024, 1024, queue),
	}
}

func (w *WebSocketServer) Run() {
	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	g.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	wsGroup := g.Group("/ws")
	{
		wsGroup.GET("/chat", w.websocketConn.WebSocketConn)
	}

	g.Run(fmt.Sprintf(":%d", w.port))
}
