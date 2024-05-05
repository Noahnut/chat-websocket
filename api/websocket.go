package api

import (
	"chat-websocket/internal"
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

	return &WebSocketServer{
		port:          port,
		websocketConn: internal.NewWebSocketConns(1024, 1024),
		ctx:           ctx,
	}
}

func (w *WebSocketServer) Run() {

	middleware := &Middleware{
		ctx: w.ctx,
	}

	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	g.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	wsGroup := g.Group("/ws")
	wsGroup.Use(middleware.AuthMiddleware)
	{
		wsGroup.GET("/chat", w.websocketConn.WebSocketConn)
	}

	g.Run(fmt.Sprintf(":%d", w.port))
}
