package internal

import (
	"chat-websocket/queue"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketConns struct {
	upgrader   websocket.Upgrader
	msgHandler *handler
}

func NewWebSocketConns(readBufferSize, writeBuffSize int, q queue.IQueue) *WebSocketConns {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBuffSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WebSocketConns{
		upgrader:   upgrader,
		msgHandler: NewHandler(q),
	}
}

func (w *WebSocketConns) WebSocketConn(c *gin.Context) {
	conn, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		resp, err := w.msgHandler.messageHandler(msg)

		if err != nil {
			log.Println(err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, resp)

		if err != nil {
			log.Println(err)
		}
	}
}
