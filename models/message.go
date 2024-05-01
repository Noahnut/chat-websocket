package models

import (
	"encoding/json"
	"net/http"
)

type MessageRequest struct {
	Type      int    `json:"type" validate:"required"`
	Target    string `json:"target" validate:"required"`
	Context   string `json:"context" validate:"required"`
	Timestamp int64  `json:"timestamp" validate:"required"`
}

type MessageResponse struct {
	Code    int    `json:"code"`
	Context string `json:"context"`
}

func (m *MessageResponse) ServerErrorResponse() []byte {
	m.Code = http.StatusInternalServerError
	m.Context = "Internal Server Error"

	resp, _ := json.Marshal(m)

	return resp
}

func (m *MessageResponse) BadRequestResponse() []byte {
	m.Code = http.StatusBadRequest
	m.Context = "Bad Request"

	resp, _ := json.Marshal(m)

	return resp
}

func (m *MessageResponse) SuccessResponse() []byte {
	m.Code = http.StatusOK
	m.Context = "Success"

	resp, _ := json.Marshal(m)

	return resp
}
