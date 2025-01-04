package statistics

import "api/internal/repository"

type GetStatisticsResponse struct {
	Year  *repository.YearStatistics  `json:"year"`
	Month *repository.MonthStatistics `json:"month"`
}
