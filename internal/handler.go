package internal

import (
	"chat-websocket/models"
	chat_protobuf "chat-websocket/models/chat-protobuf"
	"chat-websocket/streaming"
	"encoding/json"
	"log"

	"google.golang.org/protobuf/proto"
)

type handler struct {
	q streaming.IStreaming
}

func NewHandler(q streaming.IStreaming) *handler {
	return &handler{
		q: q,
	}
}

func (h *handler) messageHandler(msg []byte) []byte {
	messageReq := models.MessageRequest{}

	resp := models.MessageResponse{}

	if err := json.Unmarshal(msg, &messageReq); err != nil {
		log.Println(err)
		return resp.ServerErrorResponse()
	}

	if err := requestValidator.Struct(messageReq); err != nil {
		log.Println(err)
		return nil
	}

	proto_message := chat_protobuf.Message{
		Type:      chat_protobuf.MessageType(messageReq.Type),
		Content:   messageReq.Context,
		Timestamp: messageReq.Timestamp,
		Target:    messageReq.Target,
	}

	protoByte, err := proto.Marshal(&proto_message)

	if err != nil {
		return resp.ServerErrorResponse()
	}

	if err := h.q.Publish(h.q.GetMessageStoreSubject(), protoByte); err != nil {
		log.Println(err)
		return resp.ServerErrorResponse()
	}

	return resp.SuccessResponse()
}
