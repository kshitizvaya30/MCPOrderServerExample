package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"Supportagent/backend/internal/model"
	"Supportagent/backend/internal/service"
)

type ChatHandler struct {
	aiService *service.AIService
}

func NewChatHandler(aiService *service.AIService) *ChatHandler {
	return &ChatHandler{
		aiService: aiService,
	}
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request model.ChatRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	request.Message = strings.TrimSpace(request.Message)

	if request.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	responseMessage, err := h.aiService.GenerateResponse(
		r.Context(),
		request.Message,
	)

	if err != nil {
		http.Error(
			w,
			"failed to generate response: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := model.ChatResponse{
		Message: responseMessage,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
