package repo

import (
	"ecommerce-backend/models"
	"ecommerce-backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	DB *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetProducts(page, limit int, search, sortField, sortOrder string) ([]models.Product, int, int, error) {

	offset := (page - 1) * limit

	validFields := map[string]bool{"id": true, "name": true, "price": true}
	if !validFields[sortField] {
		sortField = "id"
	}
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	query := "SELECT * FROM products"
	args := []interface{}{}
	conditions := []string{}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+search+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)

	args = append(args, limit, offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	var products []models.Product
	err := r.DB.Select(&products, query, args...)
	if err != nil {
		log.Println("Query error:", err)
		return nil, 0, 0, err
	}

	countQuery := "SELECT COUNT(*) FROM products"
	countArgs := []interface{}{}
	if search != "" {
		countQuery += " WHERE LOWER(name) LIKE LOWER($1)"
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int
	err = r.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		log.Println("Count error:", err)
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return products, total, totalPages, nil
}

func (r *ProductRepository) GetProductByID(id int) (*models.Product, error) {
	redisKey := fmt.Sprintf("product:%d", id)

	// Check Redis
	cached, err := utils.RedisClient.Get(utils.Ctx, redisKey).Result()
	if err == nil {
		var product models.Product
		json.Unmarshal([]byte(cached), &product)
		fmt.Println("🟢 Fetched from Redis")
		return &product, nil
	}

	// Fallback to DB
	product := &models.Product{}
	err = r.DB.Get(product, "SELECT * FROM products WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	// Save to Redis
	data, _ := json.Marshal(product)
	err = utils.RedisClient.Set(utils.Ctx, redisKey, data, 0).Err() // 0 = no expiry (or use time.Hour)
	if err != nil {
		fmt.Println("🔴 Redis SET failed:", err)
	} else {
		fmt.Println("🟡 Cached in Redis:", redisKey)
	}
	fmt.Println("🔁 Writing to Redis:", redisKey, string(data))

	return product, nil
}

func (r *ProductRepository) CreateProduct(p *models.Product) error {
	query := `INSERT INTO products (name, price, image) VALUES ($1, $2, $3) RETURNING id`
	return r.DB.QueryRow(query, p.Name, p.Price, p.Image).Scan(&p.ID)
}

func (r *ProductRepository) UpdateProduct(id int, updates map[string]interface{}) (*models.Product, error) {
	setParts := []string{}
	args := []interface{}{}
	argId := 1

	for key, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s=$%d", key, argId))
		args = append(args, value)
		argId++
	}
	if len(setParts) == 0 {
		return nil, nil
	}
	args = append(args, id)

	query := fmt.Sprintf("UPDATE products SET %s WHERE id=$%d RETURNING id, name, price, image",
		strings.Join(setParts, ", "), argId)

	var product models.Product
	err := r.DB.QueryRow(query, args...).Scan(&product.ID, &product.Name, &product.Price, &product.Image)
	if err != nil {
		return nil, err
	}
	redisKey := fmt.Sprintf("product:%d", id)
	utils.RedisClient.Del(utils.Ctx, redisKey)

	return &product, nil
}

func (r *ProductRepository) DeleteProduct(id int) (bool, error) {
	res, err := r.DB.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	redisKey := fmt.Sprintf("product:%d", id)
	utils.RedisClient.Del(utils.Ctx, redisKey)

	return rowsAffected > 0, nil
}
