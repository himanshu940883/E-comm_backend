package handlers

import (
	"ecommerce-backend/middleware"
	"ecommerce-backend/models"
	"ecommerce-backend/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CartHandler struct {
	Service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{Service: service}
}

func (h *CartHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.CreateCartItem)
	r.Get("/{id}", h.GetCartForUser)
	r.Delete("/{id}", h.DeleteCartItem)
	r.Patch("/{id}", h.UpdateCartItemQuantity)
}

func (h *CartHandler) CreateCartItem(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart

	userIDFromCtx := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDFromCtx.(int)
	if !ok {
		http.Error(w, "User ID missing or invalid", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if cart.ProductID == 0 {
		http.Error(w, "user_id, product_id are required", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateCartItem(&cart, userID); err != nil {
		log.Println("CreateCartItem error:", err)
		http.Error(w, "Failed to create cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) GetCartForUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "userID parameter is required", http.StatusBadRequest)
		return
	}

	carts, err := h.Service.GetCartForUser(userID)
	if err != nil {
		log.Println("GetCartForUser error:", err)
		http.Error(w, "Failed to fetch cart items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(carts) == 0 {
		w.Write([]byte("[]"))
		return
	}
	json.NewEncoder(w).Encode(carts)
}

func (h *CartHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	deleted, err := h.Service.DeleteCartItem(id)
	if err != nil {
		log.Println("DeleteCartItem error:", err)
		http.Error(w, "Failed to delete cart item", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Cart item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	var payload struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if payload.Quantity <= 0 {
		http.Error(w, "quantity must be > 0", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.UpdateCartQuantity(id, payload.Quantity)
	if err != nil {
		log.Println("UpdateCartItemQuantity error:", err)
		http.Error(w, "Failed to update quantity", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.Error(w, "Cart item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
