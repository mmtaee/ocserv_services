package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/password"
	TokenGenerator "api/pkg/token"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserRepository struct {
	Admin AdminRepositoryInterface
	Staff StaffRepositoryInterface
	db    *gorm.DB
}

type UserRepositoryInterface interface {
	Login(context.Context, string, string, bool) (string, error)
	Logout(context.Context) error
	ChangePassword(context.Context, string, string) error
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Admin: NewAdminRepository(),
		Staff: NewStaffRepository(),
		db:    database.Connection(),
	}
}

func (r *UserRepository) Login(c context.Context, username, passwd string, rememberMe bool) (string, error) {
	var (
		user     models.User
		expireAt time.Time
	)
	err := r.db.WithContext(c).Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	passwordHash := password.Create(passwd)
	if user.Password != passwordHash {
		return "", errors.New("invalid password")
	}
	if rememberMe {
		expireAt = time.Now().Add(time.Hour * 24 * 30)
	} else {
		expireAt = time.Now().Add(time.Hour * 24)
	}
	token := models.UserToken{
		ID:        user.ID,
		Token:     TokenGenerator.Create(user.ID, expireAt),
		ExpiresAt: expireAt,
	}
	err = r.db.WithContext(c).Create(&token).Error
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (r *UserRepository) Logout(c context.Context) error {
	userID, ok := c.Value("userID").(string)
	if !ok {
		return errors.New("userID not found in context")
	}
	token, ok := c.Value("token").(string)
	if !ok {
		return errors.New("token not found in context")
	}
	return r.db.WithContext(c).
		Where("token = ? AND user_id = ? ", token, userID).
		Delete(&models.UserToken{}).Error
}

func (r *UserRepository) ChangePassword(c context.Context, oldPasswd, newPasswd string) error {
	var user models.User
	userID, ok := c.Value("userID").(string)
	if !ok {
		return errors.New("userID not found in context")
	}
	return r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return err
		}
		if user.Password != password.Create(oldPasswd) {
			return errors.New("incorrect old password")
		}
		user.Password = password.Create(newPasswd)
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})
}
