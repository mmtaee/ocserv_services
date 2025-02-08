package ocUser

import (
	"api/internal/repository"
	_ "api/internal/routes/middlewares"
	"api/pkg/utils"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/go-oc-utils/models"
	"net/http"
	"time"
)

type Controller struct {
	validator      utils.CustomValidatorInterface
	ocservUserRepo repository.OcservUserRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator:      utils.NewCustomValidator(),
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
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 pager query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Success      200  {object} OcservUsersResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users [get]
func (ctrl *Controller) Users(c echo.Context) error {
	data := utils.NewPaginationRequest()
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
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
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Success      200  {object} models.OcUser
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
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
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request body  OcservUserCreateOrUpdateRequest true "Create Ocserv User Body"
// @Success      201  {object} models.OcUser
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users [post]
func (ctrl *Controller) Create(c echo.Context) error {
	var data OcservUserCreateOrUpdateRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	user := models.OcUser{
		Group:       *data.Group,
		Username:    *data.Username,
		Password:    *data.Password,
		TrafficType: *data.TrafficType,
	}
	if data.ExpireAt != nil {
		t, err := time.Parse("2006-06-02", *data.ExpireAt)
		if err != nil {
			return utils.BadRequest(c, err)
		}
		user.ExpireAt = &t
	}
	if data.TrafficSize != nil {
		user.TrafficSize = *data.TrafficSize
	} else {
		user.TrafficSize = 0
	}
	ctx := context.WithValue(c.Request().Context(), "userID", c.Get("userID"))
	newUser, err := ctrl.ocservUserRepo.Create(ctx, &user)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, newUser)
}

// Update  Ocserv User Update
//
// @Summary      Update Ocserv User
// @Description  Update Ocserv User
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Param        request body  OcservUserCreateOrUpdateRequest true "Update Ocserv User Body"
// @Success      200  {object} models.OcUser
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users/:uid [put]
func (ctrl *Controller) Update(c echo.Context) error {
	var data OcservUserCreateOrUpdateRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	user := models.OcUser{
		Group:       *data.Group,
		Username:    *data.Username,
		Password:    *data.Password,
		TrafficType: *data.TrafficType,
	}
	if data.ExpireAt != nil {
		t, err := time.Parse("2006-06-02", *data.ExpireAt)
		if err != nil {
			return utils.BadRequest(c, err)
		}
		user.ExpireAt = &t
	}
	if data.TrafficSize != nil {
		user.TrafficSize = *data.TrafficSize
	} else {
		user.TrafficSize = 0
	}
	ctx := context.WithValue(c.Request().Context(), "userID", c.Get("userID"))
	updatedUser, err := ctrl.ocservUserRepo.Update(ctx, c.Param("uid"), &user)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, updatedUser)
}

// LockOrUnlock  Ocserv User Lock or Unlock
//
// @Summary      Lock or Unlock Ocserv User
// @Description  Lock or Unlock Ocserv User
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Param        request body  OcservUserLockRequest true "Update Ocserv User Body"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users/:uid/lock [post]
func (ctrl *Controller) LockOrUnlock(c echo.Context) error {
	var (
		data OcservUserLockRequest
		lock bool
	)
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	if data.Lock != nil {
		lock = *data.Lock
	} else {
		lock = false
	}
	ctx := context.WithValue(c.Request().Context(), "userID", c.Get("userID"))
	err := ctrl.ocservUserRepo.LockOrUnLock(ctx, c.Param("uid"), lock)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// Disconnect  Ocserv User Disconnect
//
// @Summary      Disconnect Ocserv User
// @Description  Disconnect Ocserv User
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users/:uid/disconnect [post]
func (ctrl *Controller) Disconnect(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), "userID", c.Get("userID"))
	err := ctrl.ocservUserRepo.Disconnect(ctx, c.Param("uid"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// Delete  Ocserv User Delete
//
// @Summary      Delete Ocserv User
// @Description  Delete Ocserv User
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Success      204  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users/:uid [delete]
func (ctrl *Controller) Delete(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), "userID", c.Get("userID"))
	err := ctrl.ocservUserRepo.Delete(ctx, c.Param("uid"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

// Statistics  Ocserv User Statistics
//
// @Summary      Statistics for Ocserv User
// @Description  Statistics for Ocserv User by given user UID
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Param 		 start query string false "Start date in format YYYY-MM-DD, null=time.Now()"
// @Param 		 end query string false "End date in format YYYY-MM-DD, null=time.Now().AddDate(0, 1, 0) 1 month"
// @Success      200  {object} []repository.Statistics
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users/:uid/statistics [get]
func (ctrl *Controller) Statistics(c echo.Context) error {
	var (
		err       error
		dateStart time.Time
		dateEnd   time.Time
	)
	dataStartStr := c.QueryParam("start")
	dataEndStr := c.QueryParam("end")

	if dataStartStr == "" {
		dateStart = time.Now()
		dateStart = time.Date(
			dateStart.Year(), dateStart.Month(), dateStart.Day(), 0, 0, 0, 0, dateStart.Location(),
		)
	} else {
		dateStart, err = time.Parse("2006-01-02", dataStartStr)
		if err != nil {
			return utils.BadRequest(c, err)
		}
	}
	if dataEndStr == "" {
		dateEnd = time.Now().AddDate(0, 1, 0)
		dateEnd = time.Date(
			dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 23, 59, 59, 59, dateEnd.Location(),
		)
	} else {
		dateEnd, err = time.Parse("2006-01-02", dataEndStr)
		if err != nil {
			return utils.BadRequest(c, err)
		}
	}
	stats, err := ctrl.ocservUserRepo.Statistics(c.Request().Context(), c.Param("uid"), dateStart, dateEnd)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, stats)
}

// Activities  Ocserv User Activities
//
// @Summary      Activities for Ocserv User
// @Description  Activities for Ocserv User by given user UID and Date
// @Tags         Ocserv Users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Param 		 start query string true "Start date in format YYYY-MM-DD"
// @Success      200  {object} []models.OcUserActivity
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/users/:uid/activities [get]
func (ctrl *Controller) Activities(c echo.Context) error {
	var (
		date time.Time
		err  error
	)
	if dateParam := c.QueryParam("date"); dateParam != "" {
		date, err = time.Parse("2006-01-02", dateParam)
	} else {
		date = time.Now()
	}
	if err != nil {
		return utils.BadRequest(c, err)
	}
	activities, err := ctrl.ocservUserRepo.Activity(c.Request().Context(), c.Param("uid"), date)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, activities)
}
