package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type UserRepository struct {
	QueryUsers func(ctx *echo.Context) (*[]models.User, error)
	PutUser    func(ctx *echo.Context, user *models.User) error
	GetUser    func(ctx *echo.Context, userID string) (*models.User, error)
	DeleteUser func(ctx *echo.Context, userID string) error
}

var userRepository *UserRepository

func GetUserRepository() *UserRepository {
	if userRepository == nil {
		return newUserRepository()
	}

	return userRepository
}

func newUserRepository() *UserRepository {
	repository := &UserRepository{
		QueryUsers: queryUsers,
		PutUser:    putUser,
		GetUser:    getUser,
		DeleteUser: deleteUser,
	}

	userRepository = repository
	return repository
}

func queryUsers(ctx *echo.Context) (*[]models.User, error) {
	return &[]models.User{}, nil
}

func putUser(ctx *echo.Context, user *models.User) error {
	return nil
}

func getUser(ctx *echo.Context, userID string) (*models.User, error) {
	return &models.User{}, nil
}

func deleteUser(ctx *echo.Context, userID string) error {
	return nil
}
