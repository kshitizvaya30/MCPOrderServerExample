package model

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Message string `json:"message"`
}
