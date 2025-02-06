package activity

import (
	"context"
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
	"regexp"
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

func setActivity(c context.Context, activity *models.OcUserActivity) error {
	db := database.Connection()
	return db.WithContext(c).Save(&activity).Error
}

func SetFailed(c context.Context, log string) {
	var username string

	re := regexp.MustCompile(`worker\[(.*?)\].*?(failed.*)`)
	match := re.FindStringSubmatch(log)
	if len(match) > 0 {
		username = match[1]
		log = match[2]
	} else {
		logger.Logf(logger.ERROR, "Failed to get username from log: %s", log)
		return
	}

	ocUser, err := getUser(c, username)
	if err != nil {
		return
	}

	activity := &models.OcUserActivity{
		OcUserID: ocUser.ID,
		Log:      log,
		Type:     models.Failed,
	}
	if err = setActivity(c, activity); err != nil {
		logger.Logf(logger.ERROR, "Failed to set activity: %s for user %s, %v", log, username, err)
		return
	}
}

func SetDisconnect(c context.Context, log string) {
	var (
		username string
	)
	re := regexp.MustCompile(`main\[(.*?)\].*(user\sdisconnected\s\(reason:\s[^\)]+\))`)
	if match := re.FindStringSubmatch(log); len(match) > 0 {
		username = match[1]
		log = match[2]
	} else {
		logger.Logf(logger.ERROR, "Failed to get username and reson from log: %s", log)
		return
	}

	ocUser, err := getUser(c, username)
	if err != nil {
		return
	}

	activity := &models.OcUserActivity{
		OcUserID: ocUser.ID,
		Log:      log,
		Type:     models.Disconnected,
	}
	if err = setActivity(c, activity); err != nil {
		logger.Logf(logger.ERROR, "Failed to set activity: %s for user %s, %v", log, username, err)
		return
	}
}

func SetConnect(c context.Context, log string) {
	var username string

	re := regexp.MustCompile(`main\[(.*?)\].*?(user logged in)`)
	match := re.FindStringSubmatch(log)
	if len(match) > 0 {
		username = match[1]
		log = fmt.Sprintf("%s %s", match[2], match[1])
	} else {
		logger.Logf(logger.ERROR, "Failed to get username from log: %s", log)
		return
	}

	ocUser, err := getUser(c, username)
	if err != nil {
		return
	}

	activity := &models.OcUserActivity{
		OcUserID: ocUser.ID,
		Log:      log,
		Type:     models.Connected,
	}
	if err = setActivity(c, activity); err != nil {
		logger.Logf(logger.ERROR, "Failed to set activity: %s for user %s, %v", log, username, err)
		return
	}
}
