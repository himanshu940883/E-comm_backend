package models

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID            int             `db:"id" json:"id"`
	UserID        string          `db:"user_id" json:"user_id"`
	Total         int             `db:"total" json:"total"`
	PaymentStatus string          `db:"payment_status" json:"payment_status"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	Items         json.RawMessage `db:"items" json:"items"` // raw JSON bytes, parse as needed
}
