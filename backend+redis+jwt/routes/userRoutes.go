package routes

import (
	"ecommerce-backend/handlers"
	"ecommerce-backend/repo"

	"ecommerce-backend/services"

	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router, db *sqlx.DB) {
	repo := repo.NewUserRepository(db)
	userService := services.NewUserService(repo)
	handler := handlers.NewUserHandler(userService)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", handler.GetUsers)
		r.Post("/register", handler.CreateUser)
		r.Get("/{id}", handler.GetUserByID)
		r.Patch("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
		r.Post("/login", handler.Login)
	})
}
