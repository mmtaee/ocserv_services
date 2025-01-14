package repository

import (
	"api/internal/models"
	"context"
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"gorm.io/gorm"
)

type StatisticsRepository struct {
	db *gorm.DB
}

type StatisticsRepositoryInterface interface {
	Year(c context.Context, year int) (*YearStatistics, error)
	Month(c context.Context, year, month int) (*MonthStatistics, error)
}

type YearStatistics struct {
	Year  int     `json:"year"`
	SumRx float64 `json:"sum_rx"`
	SumTx float64 `json:"sum_tx"`
}

type MonthStatistics struct {
	Year  int     `json:"year"`
	Month int     `json:"month"`
	SumRx float64 `json:"sum_rx"`
	SumTx float64 `json:"sum_tx"`
}

func NewStatisticsRepository() *StatisticsRepository {
	return &StatisticsRepository{
		db: database.Connection(),
	}
}

func (s *StatisticsRepository) Year(c context.Context, year int) (*YearStatistics, error) {
	var result YearStatistics
	err := s.db.WithContext(c).Model(&models.OcUserTrafficStatistics{}).
		Select("TO_CHAR(date, 'YYYY') AS year, SUM(rx) AS sum_rx, SUM(tx) AS sum_tx").
		Where("TO_CHAR(date, 'YYYY') = ?", fmt.Sprintf("%d", year)).
		Group("TO_CHAR(date, 'YYYY'), TO_CHAR(date, 'MM')").
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *StatisticsRepository) Month(c context.Context, year, month int) (*MonthStatistics, error) {
	var result MonthStatistics
	err := s.db.WithContext(c).Model(&models.OcUserTrafficStatistics{}).
		Select("TO_CHAR(date, 'YYYY') AS year, TO_CHAR(date, 'MM') AS month, SUM(rx) AS sum_rx, SUM(tx) AS sum_tx").
		Where("TO_CHAR(date, 'YYYY') = ? AND TO_CHAR(date, 'MM') = ?",
			fmt.Sprintf("%d", year), fmt.Sprintf("%02d", month),
		).
		Group("TO_CHAR(date, 'YYYY'), TO_CHAR(date, 'MM')").
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
