package repository

import (
	"database/sql"
	"log"

	"chat/domain"

	"github.com/pkg/errors"
)

type userRepository struct {
	DB *sql.DB
}

type UserRepository interface {
	FindAllUsers() ([]*domain.User, error)
	FindUserByID(id int) (*domain.User, error)
	FindUsreByUserID(userID string) (*domain.User, error)
	CreateUser(userID, name string) (int, error)
	DeleteUserByID(id int) error
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (ur *userRepository) FindAllUsers() ([]*domain.User, error) {
	rows, err := ur.DB.Query("SELECT id, user_id, name FROM users")
	if err != nil {
		err = errors.Wrap(err, "failed to get users")
		return nil, err
	}
	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if rows.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
			log.Printf("failed to scan: %v", err)
			continue
		}
		users = append(users, &user)
	}
	return users, nil
}

func (ur *userRepository) FindUserByID(id int) (*domain.User, error) {
	row := ur.DB.QueryRow("SELECT id, user_id, name FROM users WHERE id=?", id)
	var user domain.User
	if err := row.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return &user, nil
		}
		err = errors.Wrap(err, "failed to get user")
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) FindUsreByUserID(userID string) (*domain.User, error) {
	row := ur.DB.QueryRow("SELECT id, user_id, name FROM users WHERE user_id=?", userID)
	var user domain.User
	if err := row.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = errors.Wrap(err, "failed to get user")
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) CreateUser(userID, name string) (int, error) {
	result, err := ur.DB.Exec("INSERT INTO users(user_id, name) VALUES (?,?)", userID, name)
	if err != nil {
		err = errors.Wrap(err, "failed to insert db")
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		err = errors.Wrap(err, "failed to get new id")
		return 0, err
	}
	return int(id), nil
}

func (ur *userRepository) DeleteUserByID(id int) error {
	_, err := ur.DB.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		err = errors.Wrap(err, "failed to delete user")
		return err
	}
	return nil
}
