package service

import (
	"fmt"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

type UserService struct {
	QueryUsers func(ctx *echo.Context) (*[]models.User, error)
	CreateUser func(ctx *echo.Context, req *models.UserRequest, requestUser *models.User) (*models.User, error)
	GetUser    func(ctx *echo.Context, id string) (*models.User, error)
	UpdateUser func(ctx *echo.Context, userID string, req *models.UserRequest, existingUser *models.User) (*models.User, error)
	DeleteUser func(ctx *echo.Context, id string, existingUser *models.User) error
}

var userService *UserService

func GetUserService() *UserService {
	if userService == nil {
		return newUserService()
	}
	return userService
}

func newUserService() *UserService {
	service := &UserService{
		QueryUsers: queryUsers,
		CreateUser: createUser,
		GetUser:    getUser,
		UpdateUser: updateUser,
		DeleteUser: deleteUser,
	}

	userService = service
	return service
}

func queryUsers(ctx *echo.Context) (*[]models.User, error) {
	users, err := repository.GetUserRepository().QueryUsers(ctx)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_users", "count", len(*users))
	return users, nil
}

func getUser(ctx *echo.Context, id string) (*models.User, error) {
	userRepository := repository.GetUserRepository()

	user, err := userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, exceptions.NewErrNotFound(fmt.Sprintf("User with id %s not found", id))
	}

	logger.Action(ctx, "get_user", "userID", user.ID)
	return user, nil
}

func createUser(ctx *echo.Context, req *models.UserRequest, requestUser *models.User) (*models.User, error) {
	repository := repository.GetUserRepository()

	user := models.User{
		UserRequest: *req,
	}

	if err := repository.PutUser(ctx, &user); err != nil {
		return nil, err
	}

	logger.Action(ctx, "create_user", "userID", user.ID)
	return &user, nil
}

func updateUser(ctx *echo.Context, userID string, req *models.UserRequest, existingUser *models.User) (*models.User, error) {
	repository := repository.GetUserRepository()

	user := *existingUser
	user.UserRequest = *req

	if err := repository.PutUser(ctx, &user); err != nil {
		return nil, err
	}

	logger.Action(ctx, "update_user", "userID", user.ID)
	return &user, nil
}

func deleteUser(ctx *echo.Context, id string, existingUser *models.User) error {
	userRepository := repository.GetUserRepository()

	if err := userRepository.DeleteUser(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_user", "userID", id)
	return nil
}
