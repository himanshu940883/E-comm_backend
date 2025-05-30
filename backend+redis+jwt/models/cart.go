package models

type Cart struct {
	ID        int     `db:"id" json:"id"`
	UserID    string  `db:"user_id" json:"user_id"`
	ProductID int     `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
	Image     string  `db:"image" json:"image"`
}
