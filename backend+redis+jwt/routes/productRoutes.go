package routes

import (
	"ecommerce-backend/handlers"
	"ecommerce-backend/middleware"
	"ecommerce-backend/repo"
	"ecommerce-backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func ProductRoutes(db *sqlx.DB, redisClient *redis.Client) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuthMiddleware)

	productRepo := repo.NewProductRepository(db)
	productService := services.NewProductService(productRepo, redisClient)
	h := handlers.NewProductHandler(productService)

	r.Get("/", h.GetAllProducts)
	r.Get("/{id}", h.GetProductByID)

	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AdminOnly)

		admin.Post("/", h.CreateProduct)
		admin.Put("/{id}", h.UpdateProduct)
		admin.Delete("/{id}", h.DeleteProduct)
	})

	return r
}
