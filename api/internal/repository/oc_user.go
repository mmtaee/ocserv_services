package repository

import (
	"api/internal/models"
	"api/pkg/database"
	"api/pkg/ocserv"
	"api/pkg/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"slices"
	"sync"
)

type OcservUserRepository struct {
	db *gorm.DB
	oc *ocserv.Handler
}

type OcservUserRepositoryInterface interface {
	Users(c context.Context, page utils.RequestPagination) (*[]models.OcUser, *utils.ResponsePagination, error)
	User(c context.Context, username string) (*models.OcUser, error)
	Create(c context.Context, user *models.OcUser) error
	Update(c context.Context, uid string, user *models.OcUser) error
	LockOrUnLock(c context.Context, uid string, lock bool) error
	Disconnect(c context.Context, uid string) error
}

func NewOcservUserRepository() *OcservUserRepository {
	return &OcservUserRepository{
		db: database.Connection(),
		oc: ocserv.NewHandler(),
	}
}

func (o *OcservUserRepository) Users(c context.Context, page utils.RequestPagination) (
	*[]models.OcUser, *utils.ResponsePagination, error,
) {
	var (
		users        []models.OcUser
		totalRecords int64
		online       []string
	)
	pageResponse := utils.NewPaginationResponse()
	pageResponse.Page = page.Page
	pageResponse.PageSize = page.PageSize
	if err := o.db.WithContext(c).Model(&models.User{}).Count(&totalRecords).Error; err != nil {
		return nil, nil, err
	}
	if totalRecords == 0 {
		return &users, pageResponse, nil
	}
	pageResponse.TotalRecords = int(totalRecords)

	var wg sync.WaitGroup
	var usersErr, onlineErr error

	go func() {
		defer wg.Done()
		offset := (page.Page - 1) * page.PageSize
		order := fmt.Sprintf("%s %s", page.Order, page.Sort)
		usersErr = o.db.WithContext(c).Table("oc_users").
			Order(order).Limit(page.PageSize).Offset(offset).Scan(&users).Error
	}()

	go func() {
		ocservOnlineUsers, err := o.oc.Occtl.OnlineUsers(c)
		if err != nil {
			onlineErr = err
			return
		}
		for _, user := range ocservOnlineUsers {
			online = append(online, user.Username)
		}
	}()
	wg.Wait()

	if usersErr != nil {
		return nil, pageResponse, usersErr
	}
	if onlineErr != nil {
		return nil, pageResponse, onlineErr
	}

	for _, user := range users {
		if slices.Contains(online, user.Username) {
			user.IsOnline = true
		}
	}
	return &users, pageResponse, nil
}

func (o *OcservUserRepository) User(c context.Context, uid string) (*models.OcUser, error) {
	var user models.OcUser
	err := o.db.WithContext(c).Table("oc_users").Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	_, err = o.oc.Occtl.ShowUser(c, user.Username)
	if err != nil {
		return &user, err
	}
	user.IsOnline = true
	return &user, nil
}

func (o *OcservUserRepository) Create(c context.Context, user *models.OcUser) error {
	tx := o.db.WithContext(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table("oc_users").Create(user).Error; err != nil {
		return err
	}

	ch := make(chan error, 1)
	go func() {
		ch <- o.oc.User.CreateOrUpdateUser(c, user.Username, user.Password, user.Group)
	}()
	if err := <-ch; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (o *OcservUserRepository) Update(c context.Context, uid string, user *models.OcUser) error {
	var existing models.OcUser

	tx := o.db.WithContext(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table("oc_users").Where("uid = ?", uid).First(&existing).Error; err != nil {
		return err
	}

	existing.Group = user.Group
	existing.Username = user.Username
	existing.Password = user.Password
	existing.ExpireAt = user.ExpireAt
	existing.TrafficType = user.TrafficType
	existing.TrafficSize = user.TrafficSize

	if err := tx.Table("oc_users").Save(&existing).Error; err != nil {
		return err
	}

	ch := make(chan error, 1)
	go func() {
		ch <- o.oc.User.CreateOrUpdateUser(c, user.Username, user.Password, user.Group)
	}()
	if err := <-ch; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (o *OcservUserRepository) LockOrUnLock(c context.Context, uid string, lock bool) error {
	user := models.OcUser{}
	tx := o.db.WithContext(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table("oc_users").Where("uid = ?", uid).First(&user).Error; err != nil {
		return err
	}
	if lock {
		user.IsLocked = true
	} else {
		user.IsLocked = false
	}
	if err := tx.Table("oc_users").Save(&user).Error; err != nil {
		return err
	}
	ch := make(chan error, 1)
	go func() {
		ch <- o.oc.User.LockUnLockUser(c, user.Username, user.IsLocked)
	}()
	if err := <-ch; err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (o *OcservUserRepository) Disconnect(c context.Context, uid string) error {
	user := models.OcUser{}
	err := o.db.WithContext(c).Table("oc_users").Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return err
	}
	return o.oc.Occtl.Disconnect(c, user.Username)
}
