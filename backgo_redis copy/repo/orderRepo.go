package repo

import (
	"ecommerce-backend/models"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	DB *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateWithTx(tx *sqlx.Tx, order *models.Order) error {
	query := `
        INSERT INTO orders (user_id, total, payment_status, items)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
	err := tx.QueryRow(
		query,
		order.UserID,
		order.Total,
		order.PaymentStatus,
		order.Items,
	).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		log.Printf("Error creating order in transaction: %v", err)
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *OrderRepository) Create(order *models.Order) error {
	query := `
        INSERT INTO orders (user_id, total, payment_status, items)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
	return r.DB.QueryRow(
		query,
		order.UserID,
		order.Total,
		order.PaymentStatus,
		order.Items,
	).Scan(&order.ID, &order.CreatedAt)
}

func (r *OrderRepository) FindAll(page, limit int) ([]models.Order, error) {
	offset := (page - 1) * limit
	var orders []models.Order
	query := `SELECT * FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.DB.Select(&orders, query, limit, offset)
	return orders, err
}

func (r *OrderRepository) FindByID(id int) (models.Order, error) {
	var order models.Order
	query := `SELECT * FROM orders WHERE id = $1`
	err := r.DB.Get(&order, query, id)
	return order, err
}

func (r *OrderRepository) Update(id int, paymentStatus *string, total *int) (bool, error) {
	query := "UPDATE orders SET "
	args := []interface{}{}
	argPos := 1

	if paymentStatus != nil {
		query += "payment_status = $" + string(rune(argPos))
		args = append(args, *paymentStatus)
		argPos++
	}
	if total != nil {
		if len(args) > 0 {
			query += ", "
		}
		query += "total = $" + string(rune(argPos))
		args = append(args, *total)
		argPos++
	}

	query += " WHERE id = $" + string(rune(argPos))
	args = append(args, id)

	result, err := r.DB.Exec(query, args...)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected > 0, err
}

func (r *OrderRepository) Delete(id int) (bool, error) {
	result, err := r.DB.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	return rowsAffected > 0, err
}
