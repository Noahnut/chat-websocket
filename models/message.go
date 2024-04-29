package models

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
