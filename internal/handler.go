package internal

import (
	"chat-websocket/queue"
	"log"
)

type handler struct {
	q queue.IQueue
}

func NewHandler(q queue.IQueue) *handler {
	return &handler{
		q: q,
	}
}

func (h *handler) messageHandler(msg []byte) ([]byte, error) {

	if err := h.q.Publish(h.q.GetTextMessageSubject(), msg); err != nil {
		log.Println(err)
	}

	return msg, nil
}
