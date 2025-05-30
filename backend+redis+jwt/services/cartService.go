package services

import (
	"ecommerce-backend/models"
	"ecommerce-backend/repo"
)

type CartService struct {
	Repo *repo.CartRepository
}

func NewCartService(repo *repo.CartRepository) *CartService {
	return &CartService{Repo: repo}
}

func (s *CartService) CreateCartItem(cart *models.Cart, userIDFromJWT int) error {
	return s.Repo.Create(cart, userIDFromJWT)
}

func (s *CartService) GetCartForUser(userID string) ([]models.Cart, error) {
	return s.Repo.FindByUserID(userID)
}

func (s *CartService) DeleteCartItem(id int) (bool, error) {
	return s.Repo.DeleteByID(id)
}

func (s *CartService) UpdateCartQuantity(id int, quantity int) (bool, error) {
	if quantity <= 0 {
		return false, nil
	}
	return s.Repo.UpdateQuantity(id, quantity)
}
