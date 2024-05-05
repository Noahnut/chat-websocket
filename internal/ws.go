package internal

import (
	"chat-websocket/streaming"
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type WebSocketConns struct {
	upgrader  websocket.Upgrader
	streaming streaming.IStreaming
}

var once sync.Once
var requestValidator *validator.Validate

const streamAddr = "nats://192.168.0.109:30176"

func NewWebSocketConns(readBufferSize, writeBuffSize int) *WebSocketConns {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBuffSize,

		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WebSocketConns{
		upgrader: upgrader,
	}
}

func (w *WebSocketConns) WebSocketConn(c *gin.Context) {
	conn, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	userID, exist := c.Get("uid")

	if !exist {
		log.Println("uid not exist")
		return
	}

	ctxAny, exist := c.Get("context")

	if !exist {
		log.Println("context not exist")
		return
	}

	ctx := ctxAny.(context.Context)

	once.Do(func() {
		requestValidator = validator.New()
	})

	streaming := streaming.IStreaming(streaming.NewNATS(streamAddr))

	if err = streaming.Connect(); err != nil {
		log.Println(err)
		return
	}

	readSubjectList := []string{streaming.GetPrivateMessageSubject(userID.(string))}
	writeSubjectList := []string{streaming.GetMessageStoreSubject()}

	messageHandler := NewHandler(ctx, streaming, userID.(string), readSubjectList, writeSubjectList)

	conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	for _, subject := range readSubjectList {
		if err = streaming.Subscribe(subject, messageHandler.subscribeMessageHandler); err != nil {
			log.Println(err)
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("context done")
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			resp := messageHandler.messageHandler(msg)

			err = conn.WriteMessage(websocket.TextMessage, resp)

			if err != nil {
				log.Println(err)
			}

			conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
		}

	}
}
