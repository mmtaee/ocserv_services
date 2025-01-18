package checker

import (
	"context"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/handler/ocuser"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

func getLastDayOfMonth(t time.Time) time.Time {
	firstDayNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstDayNextMonth.Add(-24 * time.Hour)
}

func RestoreMonthlyAccounts(c context.Context) {
	db := database.Connection()
	var ocUsers []models.OcUser

	now := time.Now()
	startOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDayOfMonth := getLastDayOfMonth(now)

	err := db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("deactivated_at < ?", startOfCurrentMonth).
			Where("expire_at IS NULL AND traffic_type IN (?, ?)", models.MonthlyReceive, models.MonthlyTransmit).
			Find(&ocUsers).Error; err != nil {
			return err
		}

		if err := db.WithContext(c).Where("id IN (?)", GetIds(ocUsers)).
			Updates(map[string]interface{}{
				"is_locked":      false,
				"deactivated_at": nil,
				"expires_at":     lastDayOfMonth,
			}).Error; err != nil {
		}
		return nil
	})

	if err != nil {
		logger.Logf(logger.WARNING, "restore account failed: %v", err)
		return
	}

	ocservUser := ocuser.NewOcservUser()

	var wg sync.WaitGroup
	for _, ocUser := range ocUsers {
		user := ocUser
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ocservUser.UnLock(c, user.Username)
		}()
	}
	wg.Wait()
}
