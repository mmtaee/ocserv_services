package user

import (
	"api/internal/repository"
	"api/pkg/config"
	"api/pkg/utils"
	"api/pkg/validator"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Controller struct {
	validator validator.CustomValidatorInterface
	userRepo  *repository.UserRepository
}

func New() *Controller {
	return &Controller{
		validator: validator.NewCustomValidator(),
		userRepo:  repository.NewUserRepository(),
	}
}

// CreateSuperUser Create Superuser account
//
// @Summary      Create Superuser
// @Description  Create Superuser in initializing step
// @Tags         Site User
// @Accept       json
// @Produce      json
// @Param        secret query string true "check secret key from file 'init_secret'"
// @Param        request body  CreateAdminUserRequest true "admin user body data"
// @Success      200  {object} CreateAdminUserResponse
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/user/admin [post]
func (ctrl *Controller) CreateSuperUser(c echo.Context) error {
	var data CreateAdminUserRequest
	secret := c.QueryParam("secret")

	if secret == "" {
		return utils.BadRequest(c, errors.New("secret parameter is required"))
	}
	file := config.GetApp().InitSecretFile
	_, err := os.Stat(file)
	if err != nil {
		return utils.BadRequest(c, errors.New("secret file does not exist"))
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	if strings.TrimSpace(secret) != strings.TrimSpace(string(content)) {
		return utils.BadRequest(
			c,
			errors.New("invalid secret key or initial application preparation steps have already been completed"),
		)
	}

	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	user, err := ctrl.userRepo.Admin.CreateSuperUser(c.Request().Context(), data.Username, data.Password)
	go func() {
		err = os.Remove(file)
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return utils.BadRequest(c, err)
	}

	token, err := ctrl.userRepo.CreateToken(c.Request().Context(), user.ID, time.Now().Add(time.Hour*24*30))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, CreateAdminUserResponse{Token: token})
}

// Login Admin or Staff login
//
// @Summary      Login Admin or Staff User
// @Description  Login Admin or Staff User to get Token for request
// @Tags         Site User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request    body  LoginRequest   true "query params"
// @Success      200  {object}  LoginResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/user/login [post]
func (ctrl *Controller) Login(c echo.Context) error {
	var (
		data LoginRequest
	)
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	token, err := ctrl.userRepo.Login(c.Request().Context(), data.Username, data.Password, data.RememberMe)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// Logout Admin or Staff logout
//
// @Summary      Logout Admin or Staff User
// @Description  Logout Admin or Staff User
// @Tags         Site User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      204  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/user/logout [delete]
func (ctrl *Controller) Logout(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), "token", c.Get("token"))
	ctx = context.WithValue(ctx, "userID", c.Get("userID"))
	err := ctrl.userRepo.Logout(ctx)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ChangePassword Admin or Staff change password
//
// @Summary      ChangePassword Admin or Staff User change password
// @Description  ChangePassword Admin or Staff User change password with send old and new password
// @Tags         Site User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request    body  LoginRequest   true "query params"
// @Success      200  {object}  ChangePasswordRequest
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/user/change_password [post]
func (ctrl *Controller) ChangePassword(c echo.Context) error {
	var data ChangePasswordRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	err := ctrl.userRepo.ChangePassword(c.Request().Context(), data.OldPassword, data.NewPassword)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return nil
}
