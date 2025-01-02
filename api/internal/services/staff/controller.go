package staff

import (
	"api/internal/models"
	"api/internal/repository"
	"api/pkg/password"
	"api/pkg/utils"
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator validator.CustomValidatorInterface
	staffRepo repository.StaffRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator: validator.NewCustomValidator(),
		staffRepo: repository.NewStaffRepository(),
	}
}

// Staffs List of Staffs
//
// @Summary      Staffs
// @Description  List of Staffs
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param page query int false "Page number, starting from 1" minimum(1)
// @Param page_size query int false "Number of items per page" minimum(1) maximum(100)
// @Param order query string false "Field to order by"
// @Param sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Success      200  {object} StaffsResponse
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/staffs [get]
func (ctrl *Controller) Staffs(c echo.Context) error {
	data := utils.NewPaginationRequest()
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	staffs, meta, err := ctrl.staffRepo.Staffs(c.Request().Context(), &data)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, StaffsResponse{
		Staffs: staffs,
		Meta:   meta,
	})
}

// StaffPermission Permission of Staffs
//
// @Summary      Staff Permission
// @Description  Staff Permission
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param 		 uid path string true "User UID"
// @Success      200  {object} models.UserPermission
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/staffs/:uid/permission [get]
func (ctrl *Controller) StaffPermission(c echo.Context) error {
	permission, err := ctrl.staffRepo.Permission(c.Request().Context(), c.Param("uid"))
	if err != nil {
		return utils.BadRequest(c, err.(error))
	}
	return c.JSON(http.StatusOK, permission)
}

// CreateStaff Create Staff
//
// @Summary      Create Staff
// @Description  Create Staff with Permission
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param        request body  CreateStaffRequest true "Staff user and permission body"
// @Success      201  {object} models.User
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/staffs [post]
func (ctrl *Controller) CreateStaff(c echo.Context) error {
	var data CreateStaffRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	staff := models.User{
		Username: data.User.Username,
		Password: password.Create(data.User.Password),
	}
	permission := models.UserPermission{
		OcUser:    data.Permission.OcUser,
		OcGroup:   data.Permission.OcGroup,
		Occtl:     data.Permission.Occtl,
		Statistic: data.Permission.Statistic,
		System:    data.Permission.System,
	}
	err := ctrl.staffRepo.CreateStaff(c.Request().Context(), &staff, &permission)
	if err != nil {
		return utils.BadRequest(c, err.(error))
	}
	return c.JSON(http.StatusCreated, staff)
}

// UpdateStaffPermission Update Staff Permission
//
// @Summary      Update Staff Permission
// @Description  Update Staff Permission
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param 		 uid path string true "User UID"
// @Param        request body  UpdateStaffPermissionRequest true "Staff permission body"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/staffs/:uid/permission [patch]
func (ctrl *Controller) UpdateStaffPermission(c echo.Context) error {
	var data UpdateStaffPermissionRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	permission := models.UserPermission{
		OcUser:    data.OcUser,
		OcGroup:   data.OcGroup,
		Occtl:     data.Occtl,
		Statistic: data.Statistic,
		System:    data.System,
	}
	err := ctrl.staffRepo.UpdateStaffPermission(c.Request().Context(), c.Param("uid"), &permission)
	if err != nil {
		return utils.BadRequest(c, err.(error))
	}
	return c.JSON(http.StatusOK, nil)
}

// UpdateStaffPassword Update Staff Password
//
// @Summary      Update Staff Password
// @Description  Update Staff Password
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param 		 uid path string true "User UID"
// @Param        request body  UpdateStaffPasswordRequest true "Staff password update body"
// @Success      200  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/staffs/:uid [patch]
func (ctrl *Controller) UpdateStaffPassword(c echo.Context) error {
	var data UpdateStaffPasswordRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	err := ctrl.staffRepo.UpdateStaffPassword(c.Request().Context(), c.Param("uid"), password.Create(data.Password))
	if err != nil {
		return utils.BadRequest(c, err.(error))
	}
	return c.JSON(http.StatusOK, nil)
}

// DeleteStaff Delete Staff
//
// @Summary      Delete Staff
// @Description  Delete Staff
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param 		 uid path string true "User UID"
// @Success      204  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/staffs/:uid [delete]
func (ctrl *Controller) DeleteStaff(c echo.Context) error {
	err := ctrl.staffRepo.DeleteStaff(c.Request().Context(), c.Param("uid"))
	if err != nil {
		return utils.BadRequest(c, err.(error))
	}
	return c.JSON(http.StatusNoContent, nil)
}
