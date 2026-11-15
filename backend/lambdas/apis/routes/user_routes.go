package routes

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/models"
	"apis/responses"
	"apis/service"
	"apis/utils"

	"github.com/labstack/echo/v5"
)

type UserRoutes struct {
	QueryUsers func(c *echo.Context) error
	PostUser   func(c *echo.Context) error
	GetUser    func(c *echo.Context) error
	PutUser    func(c *echo.Context) error
	DeleteUser func(c *echo.Context) error
}

var userRoutes *UserRoutes

func getUserRoutes() *UserRoutes {
	if userRoutes == nil {
		return newUserRoutes()
	}
	return userRoutes
}

func newUserRoutes() *UserRoutes {
	routes := &UserRoutes{
		QueryUsers: queryUsers,
		PostUser:   postUser,
		GetUser:    getUser,
		PutUser:    putUser,
		DeleteUser: deleteUser,
	}

	userRoutes = routes
	return routes
}

func queryUsers(ctx *echo.Context) error {
	service := service.GetUserService()

	page, err := service.QueryUsers(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

func postUser(ctx *echo.Context) error {
	service := service.GetUserService()

	requestUser, ok := ctx.Get("requestUser").(*models.User)
	if !ok || requestUser == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get request user")
	}

	var p models.UserRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	user, err := service.CreateUser(ctx, &p, requestUser)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, user)
}

func getUser(ctx *echo.Context) error {
	userService := service.GetUserService()

	id := ctx.Param("userID")

	user, ok := ctx.Get("user").(*models.User)
	if ok && user != nil {
		return responses.Success(ctx, http.StatusOK, user)
	}

	user, err := userService.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, &exceptions.ErrNotFound{}) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return responses.Success(ctx, http.StatusOK, user)
}

func putUser(ctx *echo.Context) error {
	service := service.GetUserService()

	existing, ok := ctx.Get("user").(*models.User)
	if !ok || existing == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get existing user")
	}

	var p models.UserRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("userID")
	user, err := service.UpdateUser(ctx, id, &p, existing)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, user)
}

func deleteUser(ctx *echo.Context) error {
	service := service.GetUserService()

	user, ok := ctx.Get("user").(*models.User)
	if !ok || user == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}

	id := ctx.Param("userID")
	if err := service.DeleteUser(ctx, id, user); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
