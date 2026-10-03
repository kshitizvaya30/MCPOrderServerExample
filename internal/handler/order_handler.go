package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"Supportagent/backend/internal/repository"
	"Supportagent/backend/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

type CreateReturnRequest struct {
	Reason string `json:"reason"`
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {

	orderNumber := strings.TrimPrefix(r.URL.Path, "/orders/")

	if orderNumber == "" {
		http.Error(w, "order number is required", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrder(r.Context(), orderNumber)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) GetOrdersByUserID(
	w http.ResponseWriter,
	r *http.Request,
) {
	userIDString := strings.TrimPrefix(r.URL.Path, "/users/")

	userIDString = strings.TrimSuffix(userIDString, "/orders")

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	orders, err := h.orderService.GetOrdersByUserID(
		r.Context(),
		userID,
	)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *OrderHandler) CreateReturn(
	w http.ResponseWriter,
	r *http.Request,
) {
	orderNumber := strings.TrimPrefix(
		r.URL.Path,
		"/orders/",
	)

	orderNumber = strings.TrimSuffix(
		orderNumber,
		"/returns",
	)

	var request CreateReturnRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	returnRequest, err := h.orderService.CreateReturn(
		r.Context(),
		orderNumber,
		request.Reason,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(returnRequest); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *OrderHandler) CreateRefund(
	w http.ResponseWriter,
	r *http.Request,
) {
	orderNumber := strings.TrimPrefix(
		r.URL.Path,
		"/orders/",
	)

	orderNumber = strings.TrimSuffix(
		orderNumber,
		"/refunds",
	)

	refund, err := h.orderService.CreateRefund(
		r.Context(),
		orderNumber,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(refund); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}
