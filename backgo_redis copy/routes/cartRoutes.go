package routes

import (
	"ecommerce-backend/handlers"
	"ecommerce-backend/middleware"
	"ecommerce-backend/repo"
	"ecommerce-backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func CartRoutes(r chi.Router, db *sqlx.DB) {
	repository := repo.NewCartRepository(db)
	service := services.NewCartService(repository)
	handler := handlers.NewCartHandler(service)

	r.Route("/api/cart", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)
		handler.RegisterRoutes(r)
	})
}
