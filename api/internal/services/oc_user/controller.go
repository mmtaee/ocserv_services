package ocUser

import (
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/utils"
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator      validator.CustomValidatorInterface
	ocservUserRepo repository.OcservUserRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator:      validator.NewCustomValidator(),
		ocservUserRepo: repository.NewOcservUserRepository(),
	}
}

// Users List of Ocserv Users
//
// @Summary      List of Ocserv Users
// @Description  List of Ocserv Users with pagination
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 pager query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Success      200  {object} OcservUsersResponse
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/ocserv/users [get]
func (ctrl *Controller) Users(c echo.Context) error {
	data := utils.NewPaginationRequest()
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	users, meta, err := ctrl.ocservUserRepo.Users(c.Request().Context(), data)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, OcservUsersResponse{
		OcUsers: users,
		Meta:    meta,
	})
}

// User Retrieve Ocserv User
//
// @Summary      Retrieve Ocserv User
// @Description  Retrieve Ocserv User by given uid
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param 		 uid path string true "Ocserv User UID"
// @Success      200  {object} models.OcUser
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/ocserv/users/:uid [get]
func (ctrl *Controller) User(c echo.Context) error {
	user, err := ctrl.ocservUserRepo.User(c.Request().Context(), c.Param("uid"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, user)
}

// Create  Ocserv User Create
//
// @Summary      Create Ocserv User
// @Description  Create Ocserv User
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        request body  OcservUserCreateOrUpdateRequest true "Create Ocserv User Body"
// @Success      201  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/ocserv/users [post]
func (ctrl *Controller) Create(c echo.Context) error {
	var data OcservUserCreateOrUpdateRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	user := models.OcUser{
		Username:    data.Username,
		Password:    data.Password,
		TrafficType: data.TrafficType,
		ExpireAt:    data.ExpireAt,
	}
	if data.TrafficSize != nil {
		user.TrafficSize = *data.TrafficSize
	} else {
		user.TrafficSize = 0
	}
	err := ctrl.ocservUserRepo.Create(c.Request().Context(), &user)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

// Update  Ocserv User Update
//
// @Summary      Update Ocserv User
// @Description  Update Ocserv User
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        request body  OcservUserCreateOrUpdateRequest true "Update Ocserv User Body"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/ocserv/users/:uid [put]
func (ctrl *Controller) Update(c echo.Context) error {
	var data OcservUserCreateOrUpdateRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	user := models.OcUser{
		Group:       data.Group,
		Username:    data.Username,
		Password:    data.Password,
		ExpireAt:    data.ExpireAt,
		TrafficType: data.TrafficType,
	}
	if data.TrafficSize != nil {
		user.TrafficSize = *data.TrafficSize
	} else {
		user.TrafficSize = 0
	}
	err := ctrl.ocservUserRepo.Update(c.Request().Context(), c.Param("uid"), &user)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

// LockOrUnlock  Ocserv User Lock or Unlock
//
// @Summary      Lock or Unlock Ocserv User
// @Description  Lock or Unlock Ocserv User
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        request body  OcservUserLockRequest true "Update Ocserv User Body"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/ocserv/users/:uid/lock [post]
func (ctrl *Controller) LockOrUnlock(c echo.Context) error {
	var data OcservUserLockRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	err := ctrl.ocservUserRepo.LockOrUnLock(c.Request().Context(), c.Param("uid"), data.Lock)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (ctrl *Controller) Disconnect(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Delete(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Statistics(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Activity(c echo.Context) error {
	return nil
}
