package services

import (
	"ecommerce-backend/models"
	"ecommerce-backend/repo"
	"encoding/json"
	"fmt"
	"time"

	"context"

	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	Repo  *repo.ProductRepository
	Redis *redis.Client
	Ctx   context.Context
}

func NewProductService(repo *repo.ProductRepository, redisClient *redis.Client) *ProductService {
	return &ProductService{
		Repo:  repo,
		Redis: redisClient,
		Ctx:   context.Background(),
	}
}

func (s *ProductService) GetProducts(page, limit int, search, sortField, sortOrder string) ([]models.Product, int, int, error) {
	cacheKey := fmt.Sprintf("products:%d:%d:%s:%s:%s", page, limit, search, sortField, sortOrder)

	cached, err := s.Redis.Get(s.Ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("⚠️ Cache miss - key not found")
	} else if err != nil {
		fmt.Println("🔴 Redis GET error:", err)
	} else {
		var cachedData struct {
			Products   []models.Product `json:"products"`
			Total      int              `json:"total"`
			TotalPages int              `json:"totalPages"`
		}
		err = json.Unmarshal([]byte(cached), &cachedData)
		if err == nil {
			fmt.Println("🟢 Products fetched from Redis cache")
			return cachedData.Products, cachedData.Total, cachedData.TotalPages, nil
		} else {
			fmt.Println("🔴 Redis unmarshal error:", err)
		}
	}

	products, total, totalPages, err := s.Repo.GetProducts(page, limit, search, sortField, sortOrder)
	if err != nil {
		return nil, 0, 0, err
	}

	type cachedData struct {
		Products   []models.Product `json:"products"`
		Total      int              `json:"total"`
		TotalPages int              `json:"totalPages"`
	}

	dataToCache := cachedData{
		Products:   products,
		Total:      total,
		TotalPages: totalPages,
	}
	dataBytes, _ := json.Marshal(dataToCache)

	err = s.Redis.Set(s.Ctx, cacheKey, dataBytes, 5*time.Minute).Err()
	if err != nil {
		fmt.Println("🔴 Redis SET failed:", err)
	} else {
		fmt.Println("🟡 Cached products in Redis:", cacheKey)
	}

	return products, total, totalPages, nil
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.Repo.GetProductByID(id)
}

func (s *ProductService) CreateProduct(p *models.Product) error {
	return s.Repo.CreateProduct(p)
}

func (s *ProductService) UpdateProduct(id int, updates map[string]interface{}) (*models.Product, error) {
	return s.Repo.UpdateProduct(id, updates)
}

func (s *ProductService) DeleteProduct(id int) (bool, error) {
	return s.Repo.DeleteProduct(id)
}
