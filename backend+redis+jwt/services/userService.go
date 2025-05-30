package services

import (
	"ecommerce-backend/models"
	"ecommerce-backend/repo"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repo.UserRepository
}

func NewUserService(repo *repo.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.Repo.Create(user)
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) UpdateUser(id int, updates map[string]interface{}) (*models.User, error) {
	return s.Repo.Update(id, updates)
}

func (s *UserService) DeleteUser(id int) (bool, error) {
	return s.Repo.Delete(id)
}

func (s *UserService) Login(userID, password, role string) (*models.User, error) {
	user, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Role != role {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Password = ""
	return user, nil
}
