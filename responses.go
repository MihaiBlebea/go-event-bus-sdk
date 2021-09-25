package eventbus

import "time"

type EventHistoryResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Events  []EventResponse `json:"events,omitempty"`
}

type EventResponse struct {
	Name         string    `json:"name"`
	Url          string    `json:"url"`
	Sent         bool      `json:"sent"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type ProjectResponse struct {
	Token   string `json:"token,omitempty"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
