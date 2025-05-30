package repo

import (
	"ecommerce-backend/models"

	"github.com/jmoiron/sqlx"
)

type CartRepository struct {
	DB *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *CartRepository {
	return &CartRepository{DB: db}
}

func (r *CartRepository) Create(cart *models.Cart, jwtUserID int) error {
	var stringUserID string
	err := r.DB.Get(&stringUserID, "SELECT user_id FROM users WHERE id=$1", jwtUserID)
	if err != nil {
		return err
	}
	cart.UserID = stringUserID

	var productPrice float64
	err = r.DB.Get(&productPrice, "SELECT price FROM products WHERE id = $1", cart.ProductID)
	if err != nil {
		return err
	}

	if cart.Quantity == 0 {
		cart.Quantity = 1
	}

	var productImg string
	err = r.DB.Get(&productImg, "SELECT image FROM products WHERE id = $1", cart.ProductID)
	if err != nil {
		return err
	}

	totalPrice := float64(cart.Quantity) * productPrice
	query := `
        INSERT INTO carts (user_id, product_id, quantity, price, image)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	return r.DB.QueryRow(
		query,
		cart.UserID,
		cart.ProductID,
		cart.Quantity,
		totalPrice,
		productImg,
	).Scan(&cart.ID)
}

func (r *CartRepository) FindByUserID(userID string) ([]models.Cart, error) {
	var carts []models.Cart
	query := `SELECT * FROM carts WHERE user_id = $1`
	err := r.DB.Select(&carts, query, userID)
	return carts, err
}

func (r *CartRepository) DeleteByID(id int) (bool, error) {
	result, err := r.DB.Exec(`DELETE FROM carts WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

func (r *CartRepository) UpdateQuantity(id int, quantity int) (bool, error) {
	var productID int
	err := r.DB.Get(&productID, "SELECT product_id FROM carts WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	var productPrice float64
	err = r.DB.Get(&productPrice, "SELECT price FROM products WHERE id = $1", productID)
	if err != nil {
		return false, err
	}

	totalPrice := float64(quantity) * productPrice

	result, err := r.DB.Exec(`UPDATE carts SET quantity = $1, price = $2 WHERE id = $3`, quantity, totalPrice, id)
	if err != nil {
		return false, err
	}

	rows, _ := result.RowsAffected()
	return rows > 0, nil
}

func (r *CartRepository) DeleteByUserIDWithTx(tx *sqlx.Tx, userID string) error {
	query := `DELETE FROM carts WHERE user_id = $1`
	_, err := tx.Exec(query, userID)
	return err
}
