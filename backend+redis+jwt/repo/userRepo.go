package repo

import (
	"errors"
	"strconv"
	"strings"

	"ecommerce-backend/models"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	if user.UserID == "" || user.Password == "" || user.Role == "" {
		return errors.New("user_id, password and role are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (user_id, password, role) VALUES ($1, $2, $3) RETURNING id`
	err = r.DB.QueryRow(query, user.UserID, string(hashedPassword), user.Role).Scan(&user.ID)
	if err != nil {
		return err
	}
	user.Password = ""
	return nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Select(&users, "SELECT id, user_id, role FROM users ORDER BY id ASC")
	return users, err
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.Get(&user, "SELECT id, user_id, password, role FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(id int, updates map[string]interface{}) (*models.User, error) {
	setParts := []string{}
	args := []interface{}{}
	argId := 1

	if userID, ok := updates["user_id"].(string); ok {
		setParts = append(setParts, "user_id=$"+itoa(argId))
		args = append(args, userID)
		argId++
	}
	if password, ok := updates["password"].(string); ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		setParts = append(setParts, "password=$"+itoa(argId))
		args = append(args, string(hashedPassword))
		argId++
	}
	if role, ok := updates["role"].(string); ok {
		setParts = append(setParts, "role=$"+itoa(argId))
		args = append(args, role)
		argId++
	}

	if len(setParts) == 0 {
		return nil, errors.New("no fields to update")
	}

	args = append(args, id)
	query := "UPDATE users SET " + join(setParts) + " WHERE id=$" + itoa(argId) + " RETURNING id, user_id, role"

	var user models.User
	err := r.DB.QueryRow(query, args...).Scan(&user.ID, &user.UserID, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Delete(id int) (bool, error) {
	res, err := r.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func join(arr []string) string {
	return strings.Join(arr, ", ")
}

func (r *UserRepository) GetByUserID(userID string) (*models.User, error) {
	var user models.User
	err := r.DB.Get(&user, "SELECT id, user_id, password, role FROM users WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
