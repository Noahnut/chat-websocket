package internal

import (
	"chat-websocket/models"
	chat_protobuf "chat-websocket/models/chat-protobuf"
	"chat-websocket/streaming"
	"context"
	"encoding/json"
	"log"

	"google.golang.org/protobuf/proto"
)

type handler struct {
	ctx              context.Context
	userID           string
	readSubjectList  []string
	writeSubjectList []string
	q                streaming.IStreaming
}

func NewHandler(ctx context.Context, q streaming.IStreaming, userID string, readSubjectList, writeSubjectList []string) *handler {
	return &handler{
		userID:           userID,
		readSubjectList:  readSubjectList,
		writeSubjectList: writeSubjectList,
		q:                q,
	}
}

func (h *handler) messageHandler(msg []byte) []byte {
	req, resp := models.MessageRequest{}, models.MessageResponse{}

	if err := json.Unmarshal(msg, &req); err != nil {
		log.Println(err)
		return resp.BadRequestResponse()
	}

	if err := requestValidator.Struct(req); err != nil {
		log.Println(err)
		return resp.BadRequestResponse()
	}

	proto_message := chat_protobuf.Message{
		Type:      chat_protobuf.MessageType(req.Type),
		Content:   req.Context,
		Timestamp: req.Timestamp,
		Target:    req.Target,
	}

	protoByte, err := proto.Marshal(&proto_message)

	if err != nil {
		return resp.ServerErrorResponse()
	}

	for _, subject := range h.writeSubjectList {
		if err := h.q.Publish(subject, protoByte); err != nil {
			log.Println(err)
			return resp.ServerErrorResponse()
		}
	}

	return resp.SuccessResponse()
}
