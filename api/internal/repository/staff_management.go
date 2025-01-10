package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StaffRepository struct {
	db *gorm.DB
}

type StaffRepositoryInterface interface {
	Staffs(c context.Context, page *utils.RequestPagination) (*[]models.User, *utils.ResponsePagination, error)
	Permission(c context.Context, userUID string) (*models.UserPermission, error)
	CreateStaff(c context.Context, user *models.User, permission *models.UserPermission) error
	UpdateStaffPermission(c context.Context, userUID string, permission *models.UserPermission) error
	UpdateStaffPassword(c context.Context, userUID, password, salt string) error
	DeleteStaff(c context.Context, userUID string) error
}

func NewStaffRepository() *StaffRepository {
	return &StaffRepository{
		db: database.Connection(),
	}
}

func (s *StaffRepository) Staffs(c context.Context, page *utils.RequestPagination) (*[]models.User, *utils.ResponsePagination, error) {
	var (
		staffs       []models.User
		totalRecords int64
	)
	pageResponse := utils.NewPaginationResponse()
	pageResponse.Page = page.Page
	pageResponse.PageSize = page.PageSize
	if err := s.db.WithContext(c).Model(&models.User{}).Count(&totalRecords).Error; err != nil {
		return nil, nil, err
	}
	if totalRecords == 0 {
		return &staffs, pageResponse, nil
	}
	pageResponse.TotalRecords = int(totalRecords)
	offset := (page.Page - 1) * page.PageSize
	order := fmt.Sprintf("%s %s", page.Order, page.Sort)
	err := s.db.WithContext(c).Model(&staffs).Where("is_admin = false").
		Order(order).Limit(page.PageSize).Offset(offset).Scan(&staffs).Error
	if err != nil {
		return nil, nil, err
	}
	return &staffs, pageResponse, nil
}

func (s *StaffRepository) Permission(c context.Context, userUID string) (*models.UserPermission, error) {
	var permission models.UserPermission
	err := s.db.Joins("JOIN users ON users.id = user_permissions.user_id").
		Where("users.uid = ?", userUID).
		First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (s *StaffRepository) CreateStaff(c context.Context, staff *models.User, permission *models.UserPermission) error {
	staff.IsAdmin = false
	err := s.db.WithContext(c).Create(&staff).Error
	if err != nil {
		return err
	}
	permission.UserID = staff.ID
	return s.db.WithContext(c).Create(&permission).Error
}

func (s *StaffRepository) UpdateStaffPermission(c context.Context, userUID string, permission *models.UserPermission) error {
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		return s.db.Model(&models.UserPermission{}).
			Where("user_id = (?)", s.db.Table("users").Select("id").Where("uid = ?", userUID)).
			Updates(&permission).Error
	})
}

func (s *StaffRepository) UpdateStaffPassword(c context.Context, userUID, password, salt string) error {
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uid = ?", userUID).
			First(&user).Error; err != nil {
			return err
		}
		user.Password = password
		user.Salt = salt
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *StaffRepository) DeleteStaff(c context.Context, userUID string) error {
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Where("uid = ?", userUID).First(&user).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ? ", user.ID).Delete(&models.UserPermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&user).Error
	})
}
