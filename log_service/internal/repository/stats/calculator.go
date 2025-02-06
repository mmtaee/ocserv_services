package stats

import (
	"context"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/handler/occtl"
	"github.com/mmtaee/go-oc-utils/handler/ocuser"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
	"regexp"
	"strconv"
	"time"
)

func getUser(c context.Context, username string) (*models.OcUser, error) {
	var ocUser models.OcUser
	db := database.Connection()
	if err := db.WithContext(c).Where("username = ?", username).First(&ocUser).Error; err != nil {
		logger.Logf(logger.ERROR, "Failed to get username from database: %s", username)
		return nil, err
	}
	return &ocUser, nil
}

func saveUser(c context.Context, user *models.OcUser) error {
	db := database.Connection()
	return db.WithContext(c).Save(&user).Error
}

func createStat(c context.Context, stat *models.OcUserTrafficStatistics) error {
	db := database.Connection()
	return db.WithContext(c).Create(&stat).Error
}

func checkUserStats(user *models.OcUser) (bool, error) {
	var trafficSizeBytes = user.TrafficSize * (1 << 30)

	switch user.TrafficType {
	case models.MonthlyTransmit, models.TotallyTransmit:
		if user.Tx >= trafficSizeBytes {
			return false, nil
		}
	case models.MonthlyReceive, models.TotallyReceive:
		if user.Rx >= trafficSizeBytes {
			return false, nil
		}
	default:
		return true, nil
	}
	return true, nil
}

func lock(c context.Context, username string) error {
	oc := ocuser.NewOcservUser()
	return oc.Lock(c, username)
}

func disconnect(c context.Context, username string) error {
	oc := occtl.NewOcctl()
	return oc.Disconnect(c, username)
}

func Calculator(log string) {
	var (
		username string
		rx       int
		tx       int
		allow    bool
	)

	re := regexp.MustCompile(`main\[(.*?)\].*rx:\s*(\d+),\s*tx:\s*(\d+)`)
	match := re.FindStringSubmatch(log)
	if len(match) > 0 {
		username = match[1]
		rx, _ = strconv.Atoi(match[2])
		tx, _ = strconv.Atoi(match[3])

		c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		ocUser, err := getUser(c, username)
		if err != nil {
			return
		}

		stat := &models.OcUserTrafficStatistics{
			OcUserID: ocUser.ID,
			Rx:       rx,
			Tx:       tx,
		}

		if err = createStat(c, stat); err != nil {
			logger.Logf(logger.ERROR, "Failed to create stat for user: %s, %v", username, err)
			return
		}

		ocUser.Rx += rx
		ocUser.Tx += tx

		if err = saveUser(c, ocUser); err != nil {
			logger.Logf(logger.ERROR, "Failed to save stat for user: %s, %v", username, err)
			return
		}

		allow, err = checkUserStats(ocUser)
		if err != nil {
			logger.Logf(logger.ERROR, "Failed to check allow to user: %v", err)
			return
		}
		if !allow {
			if err = lock(c, username); err != nil {
				logger.Logf(logger.ERROR, "Failed to lock user: %v", err)
				return
			}
			if err = disconnect(c, username); err != nil {
				logger.Logf(logger.ERROR, "Failed to disconnect user: %v", err)
				return
			}
			logger.Logf(logger.INFO, "User %s locked and disconnected due to traffic limits", ocUser.Username)
		}
	} else {
		logger.Logf(logger.WARNING, "RxTxCalculator: no match found in line %s", log)
	}
}
