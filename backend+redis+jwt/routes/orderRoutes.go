package routes

import (
	"ecommerce-backend/handlers"
	"ecommerce-backend/middleware"
	"ecommerce-backend/repo"
	"ecommerce-backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func OrderRoutes(r chi.Router, db *sqlx.DB) {
	orderRepo := repo.NewOrderRepository(db)
	cartRepo := repo.NewCartRepository(db)
	service := services.NewOrderService(orderRepo, cartRepo)
	handler := handlers.NewOrderHandler(service)

	r.Route("/api/orders", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)
		handler.RegisterRoutes(r)
	})
}
