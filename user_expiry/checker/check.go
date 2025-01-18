package checker

import (
	"context"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/handler/occtl"
	"github.com/mmtaee/go-oc-utils/handler/ocuser"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

func CheckExpiry(c context.Context) {
	db := database.Connection()
	var ocUsers []models.OcUser

	err := db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("expire_at >= ? AND deactivated_at IS NULL AND traffic_type IN (?, ?)",
				time.Now(), models.MonthlyReceive, models.MonthlyTransmit).
			Find(&ocUsers).Error; err != nil {
			return err
		}

		if err := db.WithContext(c).
			Model(&models.OcUser{}).
			Where("id IN ?", GetIds(ocUsers)).
			Updates(map[string]interface{}{
				"is_locked":      true,
				"deactivated_at": time.Now(),
				"expires_at":     nil,
			}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Logf(logger.WARNING, "check expiry failed: %v", err)
		return
	}

	oc := occtl.NewOcctl()
	ocservUser := ocuser.NewOcservUser()
	var wg sync.WaitGroup

	for _, ocUser := range ocUsers {
		user := ocUser
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = oc.Disconnect(c, user.Username)
			_ = ocservUser.Lock(c, user.Username)
		}()
	}
	wg.Wait()
}
