package occtl

import (
	"api/internal/repository"
	_ "api/internal/routes/middlewares"
	"api/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator utils.CustomValidatorInterface
	os        repository.OcctlRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator: utils.NewCustomValidator(),
		os:        repository.NewOcctlRepository(),
	}
}

// Reload  		 server config reload
//
// @Summary      Reload Server
// @Description  Reload Server Configuration
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      202  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/reload [post]
func (ctrl *Controller) Reload(c echo.Context) error {
	err := ctrl.os.Reload(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}

// OnlineUsers   Online Users
//
// @Summary      Online User List
// @Description  Show Online User List in Server
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      200  {array} occtl.OcUser
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/online [get]
func (ctrl *Controller) OnlineUsers(c echo.Context) error {
	users, err := ctrl.os.OnlineUsers(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, users)
}

// Disconnect    Close Session Of Connected User
//
// @Summary      Disconnect User From Server
// @Description  Disconnect User From Server and Close Session
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 username path string true "Ocserv User Username"
// @Success      202  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/disconnect/:username [post]
func (ctrl *Controller) Disconnect(c echo.Context) error {
	err := ctrl.os.Disconnect(c.Request().Context(), c.Param("username"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}

// ShowIPBans    IP Bans
//
// @Summary      List Of IP Bans
// @Description  List Of IP Bans
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 points query bool false "Show IP bans with Points"
// @Success      200  {array} occtl.IPBan
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/ip_bans [get]
func (ctrl *Controller) ShowIPBans(c echo.Context) error {
	bans, err := ctrl.os.ShowIPBans(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, bans)
}

// ShowIPBansPoint    IP Bans Point
//
// @Summary      List Of IP Bans include points
// @Description  List Of IP Bans include points
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 points query bool false "Show IP bans with Points"
// @Success      200  {array} occtl.IPBanPoints
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/ip_bans/point [get]
func (ctrl *Controller) ShowIPBansPoint(c echo.Context) error {
	bans, err := ctrl.os.ShowIPBansPoint(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, bans)
}

// UnBanIP       Unban IP
//
// @Summary      Unban IP
// @Description  Remove IP From Baned List
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param		 request body  UnBanIPRequest true "IP Address"
// @Success      202  {object} nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/unban [post]
func (ctrl *Controller) UnBanIP(c echo.Context) error {
	var data UnBanIPRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	err := ctrl.os.UnBanIP(c.Request().Context(), data.IP)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}

// ShowStatus    Show Status
//
// @Summary      Show Status
// @Description  Show Status Of Server State
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      200  {object} ShowStatusResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/status [get]
func (ctrl *Controller) ShowStatus(c echo.Context) error {
	status := ctrl.os.ShowStatus(c.Request().Context())
	return c.JSON(http.StatusOK, ShowStatusResponse{
		Status: status,
	})
}

// ShowIRoutes   Show IP Routes
//
// @Summary      Show IP Routes
// @Description  Show IP Routes
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      200  {array} occtl.IRoute
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/iroutes [get]
func (ctrl *Controller) ShowIRoutes(c echo.Context) error {
	routes, err := ctrl.os.ShowIRoutes(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, routes)
}

// ShowUser      Show User Status
//
// @Summary      Show User Status
// @Description  Show User Status
// @Tags         Occtl
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 username path string true "Ocserv User Username"
// @Success      200  {array} occtl.OcUser
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/occtl/users/:username [get]
func (ctrl *Controller) ShowUser(c echo.Context) error {
	user, err := ctrl.os.ShowUser(c.Request().Context(), c.Param("username"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, user)
}
