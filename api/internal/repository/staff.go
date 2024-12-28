package repository

import (
	"api/pkg/database"
	"gorm.io/gorm"
)

type StaffRepository struct {
	db *gorm.DB
}

type StaffRepositoryInterface interface{}

func NewStaffRepository() *StaffRepository {
	return &StaffRepository{
		db: database.Connection(),
	}
}
