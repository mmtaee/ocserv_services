package repository

import (
	"api/pkg/event"
	"api/pkg/utils"
	"context"
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/handler/occtl"
	"github.com/mmtaee/go-oc-utils/handler/ocuser"
	"github.com/mmtaee/go-oc-utils/models"
	"gorm.io/gorm"
	"slices"
	"sync"
	"time"
)

type OcservUserRepository struct {
	db          *gorm.DB
	ocUser      ocuser.OcservUserInterface
	occtl       occtl.OcInterface
	WorkerEvent *event.WorkerEvent
}

type OcservUserRepositoryInterface interface {
	Users(c context.Context, page utils.RequestPagination) (*[]models.OcUser, *utils.ResponsePagination, error)
	User(c context.Context, username string) (*models.OcUser, error)
	Create(c context.Context, user *models.OcUser) (*models.OcUser, error)
	Update(c context.Context, uid string, user *models.OcUser) (*models.OcUser, error)
	LockOrUnLock(c context.Context, uid string, lock bool) error
	Disconnect(c context.Context, uid string) error
	Delete(c context.Context, uid string) error
	Statistics(c context.Context, uid string, startDate, endDate time.Time) (*[]Statistics, error)
	Activity(c context.Context, uid string, date time.Time) (*[]models.OcUserActivity, error)
}

type Statistics struct {
	CreatedAt string  `json:"created_at"`
	SumRx     float64 `json:"sum_rx"`
	SumTx     float64 `json:"sum_tx"`
}

func NewOcservUserRepository() *OcservUserRepository {
	return &OcservUserRepository{
		db:          database.Connection(),
		ocUser:      ocuser.NewOcservUser(),
		occtl:       occtl.NewOcctl(),
		WorkerEvent: event.GetWorker(),
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
		onlineUsers, err := o.occtl.OnlineUsers(c)
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
	_, err = o.occtl.ShowUser(c, user.Username)
	if err != nil {
		return &user, err
	}
	user.IsOnline = true
	return &user, nil
}

func (o *OcservUserRepository) Create(c context.Context, user *models.OcUser) (*models.OcUser, error) {
	tx := o.db.WithContext(c).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Table("oc_users").Create(user).Error; err != nil {
		return nil, err
	}
	if err := o.ocUser.Create(c, user.Username, user.Password, user.Group); err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: "create_oc_user",
		ModelName: "oc_user",
		ModelUID:  user.UID,
		UserUID:   c.Value("userID").(string),
		OldState:  nil,
		NewState:  user,
	})

	return user, nil
}

func (o *OcservUserRepository) Update(c context.Context, uid string, user *models.OcUser) (*models.OcUser, error) {
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
		return nil, err
	}

	oldState := existing

	existing.Group = user.Group
	existing.Username = user.Username
	existing.Password = user.Password
	existing.ExpireAt = user.ExpireAt
	existing.TrafficType = user.TrafficType
	existing.TrafficSize = user.TrafficSize

	if err := tx.Table("oc_users").Save(&existing).Error; err != nil {
		return nil, err
	}

	if err := o.ocUser.Update(c, user.Username, user.Password, user.Group); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: "update_oc_user",
		ModelName: "oc_user",
		ModelUID:  uid,
		UserUID:   c.Value("userID").(string),
		OldState:  oldState,
		NewState:  existing,
	})

	return &existing, nil
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

	if user.IsLocked {
		if err := o.ocUser.Lock(c, user.Username); err != nil {
			return err
		}
	} else {
		if err := o.ocUser.UnLock(c, user.Username); err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	var eventType, newState, oldState string

	if lock {
		eventType = "lock_oc_user"
		oldState = "unlock"
		newState = "lock"
	} else {
		eventType = "unlock_oc_user"
		oldState = "lock"
		newState = "unlock"
	}
	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: eventType,
		ModelName: "oc_user",
		ModelUID:  uid,
		UserUID:   c.Value("userID").(string),
		OldState:  oldState,
		NewState:  newState,
	})

	return nil
}

func (o *OcservUserRepository) Disconnect(c context.Context, uid string) error {
	user := models.OcUser{}
	err := o.db.WithContext(c).Table("oc_users").Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return err
	}
	err = o.occtl.Disconnect(c, user.Username)
	if err != nil {
		return err
	}
	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: "disconnect_oc_user",
		ModelName: "oc_user",
		ModelUID:  uid,
		UserUID:   c.Value("userID").(string),
		OldState:  nil,
		NewState:  nil,
	})
	return nil
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
	if err := o.ocUser.Delete(c, user.Username); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: "delete_oc_user",
		ModelName: "oc_user",
		ModelUID:  uid,
		UserUID:   c.Value("userID").(string),
	})

	return nil
}

func (o *OcservUserRepository) Statistics(c context.Context, uid string, startDate, endDate time.Time) (
	*[]Statistics, error,
) {
	var results []Statistics
	err := o.db.WithContext(c).
		Table("oc_user_traffic_statistics").
		Joins("JOIN oc_users ON oc_users.id = oc_user_traffic_statistics.oc_user_id").
		Where("oc_users.uid = ? AND oc_user_traffic_statistics.created_at BETWEEN ? AND ?", uid, startDate, endDate).
		Select(
			"oc_user_traffic_statistics.created_at, " +
				"SUM(oc_user_traffic_statistics.rx) as sum_rx, " +
				"SUM(oc_user_traffic_statistics.tx) as sum_tx",
		).
		Group("oc_user_traffic_statistics.created_at").
		Order("oc_user_traffic_statistics.created_at").
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
