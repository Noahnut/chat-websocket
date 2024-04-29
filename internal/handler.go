package internal

import (
	"chat-websocket/models"
	chat_protobuf "chat-websocket/models/chat-protobuf"
	"chat-websocket/queue"
	"encoding/json"
	"log"

	"google.golang.org/protobuf/proto"
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
	messageReq := models.MessageRequest{}
	if err := json.Unmarshal(msg, &messageReq); err != nil {
		log.Println(err)
		return nil, err
	}

	if err := requestValidator.Struct(messageReq); err != nil {
		log.Println(err)
		return nil, err
	}

	proto_message := chat_protobuf.Message{
		Type:      chat_protobuf.MessageType(messageReq.Type),
		Content:   messageReq.Context,
		Timestamp: messageReq.Timestamp,
		Target:    messageReq.Target,
	}

	protoByte, err := proto.Marshal(&proto_message)

	if err != nil {
		return nil, err
	}

	if err := h.q.Publish(h.q.GetTextMessageSubject(), protoByte); err != nil {
		log.Println(err)
	}

	return msg, nil
}
