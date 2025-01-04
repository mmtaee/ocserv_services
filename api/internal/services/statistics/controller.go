package statistics

import (
	"api/internal/repository"
	"api/pkg/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	statistics repository.StatisticsRepositoryInterface
}

func New() *Controller {
	return &Controller{
		statistics: repository.NewStatisticsRepository(),
	}
}

// Statistics Year and Month Statistics
//
// @Summary      Year and Month Statistics
// @Description  Year and Month Statistics by given year and month as int
// @Tags         Statistics
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 year query int true "Year"
// @Param 		 month query int true "Month number"
// @Success      200  {object}  GetStatisticsResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/statistics [get]
func (ctrl *Controller) Statistics(c echo.Context) error {
	var (
		year, month int
		err         error
	)
	monthStr := c.QueryParam("month")
	if monthStr == "" {
		month = int(time.Now().Month())
	} else {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			return utils.BadRequest(c, errors.New("invalid month"))
		}
	}
	yearStr := c.QueryParam("year")
	if yearStr == "" {
		year = time.Now().Year()
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return utils.BadRequest(c, errors.New("invalid year"))
		}
	}
	statsMonth, err := ctrl.statistics.Month(c.Request().Context(), year, month)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	statsYear, err := ctrl.statistics.Year(c.Request().Context(), year)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, GetStatisticsResponse{
		Year:  statsYear,
		Month: statsMonth,
	})
}
