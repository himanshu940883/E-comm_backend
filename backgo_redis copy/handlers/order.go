package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"ecommerce-backend/middleware"
	"ecommerce-backend/services"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.CreateOrder)
	r.Get("/", h.GetOrders)
	r.Get("/{id}", h.GetOrderByID)
	r.Patch("/{id}", h.UpdateOrder)
	r.Delete("/{id}", h.DeleteOrder)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// var payload struct {
	// 	UserID string `json:"user_id"`
	// }
	// if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
	// 	http.Error(w, "invalid request body", http.StatusBadRequest)
	// 	return
	// }

	// if payload.UserID == "" {
	// 	http.Error(w, "user_id is required", http.StatusBadRequest)
	// 	return
	// }

	jwtUserIDFromCtx := r.Context().Value(middleware.UserIDKey)
	jwtUserID, ok := jwtUserIDFromCtx.(int)
	if !ok {
		http.Error(w, "User ID missing or invalid", http.StatusUnauthorized)
		return
	}

	// 2. Lookup string user ID from DB using UserRepo or DB access inside OrderService (better)
	var stringUserID string
	err := h.service.OrderRepo.DB.Get(&stringUserID, "SELECT user_id FROM users WHERE id=$1", jwtUserID)
	if err != nil {
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	order, err := h.service.CreateOrderFromCart(stringUserID, "success")
	if err != nil {
		log.Printf("CreateOrderFromCart error: %v", err)
		if err.Error() == "cart is empty" {
			http.Error(w, "cart is empty", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	log.Printf("Order created: %+v", order)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePaginationParams(r)

	orders, err := h.service.GetOrders(page, limit)
	if err != nil {
		http.Error(w, "failed to fetch orders", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch order", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload struct {
		PaymentStatus *string `json:"payment_status"`
		Total         *int    `json:"total"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if payload.PaymentStatus == nil && payload.Total == nil {
		http.Error(w, "nothing to update", http.StatusBadRequest)
		return
	}

	ok, err := h.service.UpdateOrder(id, payload.PaymentStatus, payload.Total)
	if err != nil {
		http.Error(w, "failed to update order", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := h.service.DeleteOrder(id)
	if err != nil {
		http.Error(w, "failed to delete order", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDParam(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.Atoi(idStr)
}

func parsePaginationParams(r *http.Request) (page, limit int) {
	page, limit = 1, 10
	if p := r.URL.Query().Get("page"); p != "" {
		if i, err := strconv.Atoi(p); err == nil && i > 0 {
			page = i
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if i, err := strconv.Atoi(l); err == nil && i > 0 {
			limit = i
		}
	}
	return
}
