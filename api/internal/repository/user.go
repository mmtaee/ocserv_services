package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/password"
	TokenGenerator "api/pkg/token"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	Admin AdminRepositoryInterface
	Staff StaffRepositoryInterface
	db    *gorm.DB
}

type UserRepositoryInterface interface {
	Login(context.Context) (string, error)
	Logout(context.Context) error
	ChangePassword(context.Context) error
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Admin: NewAdminRepository(),
		Staff: NewStaffRepository(),
		db:    database.Connection(),
	}
}

func (r *UserRepository) Login(c context.Context) (string, error) {
	var (
		user     models.User
		expireAt time.Time
	)
	ch := make(chan error, 1)
	go func() {
		err := r.db.WithContext(c).Where("username = ?", c.Value("username").(string)).First(&user).Error
		if err != nil {
			ch <- err
		}
		passwordString := c.Value("password").(string)
		passwordHash := password.Create(passwordString)
		if user.Password != passwordHash {
			ch <- errors.New("invalid password")
		}
		ch <- nil
	}()
	if err := <-ch; err != nil {
		return "", err
	}
	if rememberMe := c.Value("rememberMe").(bool); rememberMe {
		expireAt = time.Now().Add(time.Hour * 24 * 30)
	} else {
		expireAt = time.Now().Add(time.Hour * 24)
	}
	token := models.UserToken{
		ID:        user.ID,
		Token:     TokenGenerator.Create(user.ID, expireAt),
		ExpiresAt: expireAt,
	}
	go func() {
		ch <- r.db.WithContext(c).Create(&token).Error
	}()
	if err := <-ch; err != nil {
		return "", err
	}
	return token.Token, nil
}

func (r *UserRepository) Logout(c context.Context) error {
	ch := make(chan error, 1)
	go func() {
		ch <- r.db.WithContext(c).
			Where("token = ? AND user_id = ? ", c.Value("token").(string), c.Value("userID").(string)).
			Delete(&models.UserToken{}).Error
	}()
	return <-ch
}

func (r *UserRepository) ChangePassword(c context.Context) error {
	// TODO: first get last password compare with request, then update.
	return nil
}
