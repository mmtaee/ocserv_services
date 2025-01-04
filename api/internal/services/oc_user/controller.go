package ocUser

import (
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
// @Success      200  {object} OcUsersResponse
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
	return c.JSON(http.StatusOK, OcUsersResponse{
		OcUsers: users,
		Meta:    meta,
	})
}

func (ctrl *Controller) User(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Create(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Update(c echo.Context) error {
	return nil
}

func (ctrl *Controller) Lock(c echo.Context) error {
	// TODO: get lock true or false from body
	return nil
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
