package services

import (
	"ecommerce-backend/models"
	"ecommerce-backend/repo"
	"encoding/json"
	"errors"
)

type OrderService struct {
	OrderRepo *repo.OrderRepository
	CartRepo  *repo.CartRepository
}

func NewOrderService(orderRepo *repo.OrderRepository, cartRepo *repo.CartRepository) *OrderService {
	return &OrderService{OrderRepo: orderRepo, CartRepo: cartRepo}
}

func (s *OrderService) CreateOrderFromCart(userID, paymentStatus string) (*models.Order, error) {
	cartItems, err := s.CartRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	var orderItems []struct {
		ProductID int     `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
		Image     string  `json:"image"`
	}
	calculatedTotal := 0.0
	for _, item := range cartItems {
		if item.Quantity <= 0 {
			return nil, errors.New("invalid quantity in cart")
		}
		calculatedTotal += item.Price * float64(item.Quantity)
		orderItems = append(orderItems, struct {
			ProductID int     `json:"product_id"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
			Image     string  `json:"image"`
		}{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Image:     item.Image,
		})
	}

	total := int(calculatedTotal)

	itemsJSON, err := json.Marshal(orderItems)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		UserID:        userID,
		Total:         total,
		PaymentStatus: paymentStatus,
		Items:         itemsJSON,
	}

	// Use transaction to create order and clear cart
	tx, err := s.OrderRepo.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.OrderRepo.CreateWithTx(tx, order); err != nil {
		return nil, err
	}

	// Clear cart
	if err := s.CartRepo.DeleteByUserIDWithTx(tx, userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrders(page, limit int) ([]models.Order, error) {
	return s.OrderRepo.FindAll(page, limit)
}

func (s *OrderService) GetOrderByID(id int) (models.Order, error) {
	return s.OrderRepo.FindByID(id)
}

func (s *OrderService) UpdateOrder(id int, paymentStatus *string, total *int) (bool, error) {
	return s.OrderRepo.Update(id, paymentStatus, total)
}

func (s *OrderService) DeleteOrder(id int) (bool, error) {
	return s.OrderRepo.Delete(id)
}
