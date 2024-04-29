package internal

import (
	"chat-websocket/queue"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type WebSocketConns struct {
	upgrader   websocket.Upgrader
	msgHandler *handler
}

var once sync.Once
var requestValidator *validator.Validate

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

	once.Do(func() {
		requestValidator = validator.New()
	})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
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
