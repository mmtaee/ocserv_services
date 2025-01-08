package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/password"
	"context"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

type AdminRepositoryInterface interface {
	CreateSuperUser(c context.Context, username, passwd string) (*models.User, error)
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		db: database.Connection(),
	}
}

func (a *AdminRepository) CreateSuperUser(c context.Context, username, passwd string) (*models.User, error) {
	pass := password.NewPassword(passwd)
	user := models.User{
		Username: username,
		Password: pass.Hash,
		Salt:     pass.Salt,
		IsAdmin:  true,
	}

	tx := a.db.WithContext(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(&user).Error
	if err != nil {

		return nil, err
	}

	if err = tx.WithContext(c).Create(&models.UserPermission{
		UserID:    user.ID,
		OcUser:    true,
		OcGroup:   true,
		Statistic: true,
		Occtl:     true,
		System:    true,
	}).Error; err != nil {
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}
	return &user, nil
}
