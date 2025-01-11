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
	"time"
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
	Delete(c context.Context, uid string) error
	Statistics(c context.Context, uid string, startDate, endDate time.Time) (*[]Statistics, error)
	Activity(c context.Context, uid string, date time.Time) (*[]models.OcUserActivity, error)
	BuildDelete(c context.Context, uid *[]string) error
}

type Statistics struct {
	Date  string  `json:"date"`
	SumRx float64 `json:"sum_rx"`
	SumTx float64 `json:"sum_tx"`
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
		users             []models.OcUser
		totalRecords      int64
		ocservOnlineUsers []string
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

	ch := make(chan error, 1)
	go func() {
		onlineUsers, err := o.oc.Occtl.OnlineUsers(c)
		if err != nil {
			ch <- err
			return
		}
		for _, u := range *onlineUsers {
			ocservOnlineUsers = append(ocservOnlineUsers, u.Username)
		}
		ch <- nil
	}()
	if err := <-ch; err != nil {
		return nil, pageResponse, err
	}

	offset := (page.Page - 1) * page.PageSize
	order := fmt.Sprintf("%s %s", page.Order, page.Sort)
	if err := o.db.WithContext(c).Table("oc_users").
		Order(order).Limit(page.PageSize).Offset(offset).Scan(&users).Error; err != nil {
		return nil, pageResponse, err
	}

	var wg sync.WaitGroup
	wg.Add(len(users))
	for i := range users {
		user := users[i]
		go func(user *models.OcUser) {
			defer wg.Done()
			if slices.Contains(ocservOnlineUsers, user.Username) {
				user.IsOnline = true
			}
		}(&user)
	}
	wg.Wait()

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
	if err := o.oc.User.CreateOrUpdateUser(c, user.Username, user.Password, user.Group); err != nil {
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

	if err := o.oc.User.CreateOrUpdateUser(c, user.Username, user.Password, user.Group); err != nil {
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

	if err := o.oc.User.LockUnLockUser(c, user.Username, user.IsLocked); err != nil {
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
func (o *OcservUserRepository) Delete(c context.Context, uid string) error {
	var user models.OcUser
	tx := o.db.WithContext(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Where("uid = ?", uid).First(&user).Error; err != nil {
		return err
	}
	if err := tx.Table("oc_users").Where("uid = ?", uid).Delete(&user).Error; err != nil {
		return err
	}
	if err := o.oc.User.DeleteUser(c, user.Username); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (o *OcservUserRepository) Statistics(c context.Context, uid string, startDate, endDate time.Time) (
	*[]Statistics, error,
) {
	var results []Statistics
	err := o.db.WithContext(c).
		Table("oc_user_traffic_statistics").
		Joins("JOIN oc_users ON oc_users.id = oc_user_traffic_statistics.oc_user_id").
		Where("oc_users.uid = ? AND date BETWEEN ? AND ?", uid, startDate, endDate).
		Select(
			"oc_user_traffic_statistics.date, " +
				"SUM(oc_user_traffic_statistics.rx) as sum_rx, " +
				"SUM(oc_user_traffic_statistics.tx) as sum_tx",
		).
		Group("oc_user_traffic_statistics.date").
		Order("oc_user_traffic_statistics.date").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (o *OcservUserRepository) Activity(c context.Context, uid string, date time.Time) (*[]models.OcUserActivity, error) {
	var activities []models.OcUserActivity
	startOfDay := date.Format("2006-01-02") + " 00:00:00"
	endOfDay := date.Format("2006-01-02") + " 23:59:59"
	err := o.db.WithContext(c).Table("oc_user_activities").
		Joins("JOIN oc_users ON oc_users.id = oc_user_activities.oc_user_id").
		Where("oc_users.uid = ? AND oc_user_activities.created_at BETWEEN ? AND ?", uid, startOfDay, endOfDay).
		Order("oc_user_activities.created_at").
		Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return &activities, nil
}

func (o *OcservUserRepository) BuildDelete(c context.Context, uid *[]string) error {
	return nil
}
