package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"context"
	"fmt"
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
		Select("?, SUM(rx) AS sum_rx, SUM(tx) AS sum_tx", year).
		Where("strftime('%Y', date) = ?", fmt.Sprintf("%d", year)).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *StatisticsRepository) Month(c context.Context, year, month int) (*MonthStatistics, error) {
	var result MonthStatistics
	err := s.db.WithContext(c).Model(&models.OcUserTrafficStatistics{}).
		Select("?, ? AS month, SUM(rx) AS sum_rx, SUM(tx) AS sum_tx", year, month).
		Where("strftime('%Y', date) = ? AND strftime('%m', date) = ?",
			fmt.Sprintf("%d", year), fmt.Sprintf("%02d", month),
		).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
