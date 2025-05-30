package main

import (
	"ecommerce-backend/db"
	"ecommerce-backend/routes"
	"ecommerce-backend/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db.Connect()

	r := chi.NewRouter()

	redisClient := utils.InitRedis()
	log.Println("✅ Connected to Redis!")

	routes.CartRoutes(r, db.DB)
	routes.OrderRoutes(r, db.DB)
	r.Mount("/api/products", routes.ProductRoutes(db.DB, redisClient))
	routes.UserRoutes(r, db.DB)

	log.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
