package user

import (
	"api/internal/repository"
	"api/pkg/utils"
	"api/pkg/validator"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator validator.CustomValidatorInterface
	userRepo  repository.UserRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator: validator.NewCustomValidator(),
		userRepo:  repository.NewUserRepository(),
	}
}

// Login Admin or Staff login
//
// @Summary      Login Admin or Staff User
// @Description  Login Admin or Staff User to get Token for request
// @Tags         Site User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      401 {object} middlewares.Unauthorized
// @Param        request    body  LoginRequest   true "query params"
// @Success      200  {object}  LoginResponse
// @Failure      400 {object} utils.ErrorResponse
// @Router       /services/v1/user/login [post]
func (ctrl *Controller) Login(c echo.Context) error {
	var data LoginRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	token, err := ctrl.userRepo.Login(c.Request().Context(), data.Username, data.Password, data.RememberMe)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}

// Logout Admin or Staff logout
//
// @Summary      Logout Admin or Staff User
// @Description  Logout Admin or Staff User
// @Tags         Site User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      401 {object} middlewares.Unauthorized
// @Success      204  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /services/v1/user/logout [delete]
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
// @Failure      401 {object} middlewares.Unauthorized
// @Param        request    body  LoginRequest   true "query params"
// @Success      200  {object}  ChangePasswordRequest
// @Failure      400 {object} utils.ErrorResponse
// @Router       /services/v1/user/change_password [post]
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
