package internal

import (
	"chat-websocket/models"
	chat_protobuf "chat-websocket/models/chat-protobuf"
	"chat-websocket/streaming"
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type handler struct {
	ctx              context.Context
	wsConn           *websocket.Conn
	userID           string
	readSubjectList  []string
	writeSubjectList []string
	q                streaming.IStreaming
}

func NewHandler(ctx context.Context, conn *websocket.Conn, q streaming.IStreaming, userID string, readSubjectList, writeSubjectList []string) *handler {
	return &handler{
		ctx:              ctx,
		wsConn:           conn,
		userID:           userID,
		readSubjectList:  readSubjectList,
		writeSubjectList: writeSubjectList,
		q:                q,
	}
}

func (h *handler) subscribeMessageHandler(msg []byte) {
	if err := h.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println(err)
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
