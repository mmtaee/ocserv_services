package activity

import (
	"context"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
	"regexp"
	"strings"
	"time"
)

type Activity struct {
	Ch chan<- string
}

var activityChan chan string

func NewActivityService() *Activity {
	activityChan = make(chan string, 100)
	go activities()

	return &Activity{Ch: activityChan}
}

func activities() {
	actionMap := map[string]func(string){
		"disconnected":          func(text string) { go setDisconnect(text) },
		"failed authentication": func(text string) { go setFailed(text) },
		"user logged in":        func(text string) { go setConnect(text) },
	}

	for msg := range activityChan {
		for keyword, action := range actionMap {
			if strings.Contains(msg, keyword) {
				go action(msg)
			}
		}
	}
}

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
	return db.WithContext(c).Save(activity).Error
}

func setFailed(log string) {
	var username string

	re := regexp.MustCompile(`worker\[(.*?)\]`)
	match := re.FindStringSubmatch(log)
	if len(match) > 0 {
		username = match[1]
	} else {
		logger.Logf(logger.ERROR, "Failed to get username from log: %s", log)
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		logger.Logf(logger.ERROR, "Failed to set activity: %s for user %s", log, username)
		return
	}
	logger.InfoF("Successfully set activity: %s for user %s", log, username)
}

func setDisconnect(log string) {
	var (
		username string
	)
	re := regexp.MustCompile(`main\[(.*?)\]`)
	if match := re.FindStringSubmatch(log); len(match) > 0 {
		username = match[1]
	} else {
		logger.Logf(logger.ERROR, "Failed to get username and reson from log: %s", log)
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		logger.Logf(logger.ERROR, "Failed to set activity: %s for user %s", log, username)
		return
	}
	logger.InfoF("Successfully set activity: %s for user %s", log, username)
}

func setConnect(log string) {
	var username string

	re := regexp.MustCompile(`main\[(.*?)\]`)
	match := re.FindStringSubmatch(log)
	if len(match) > 0 {
		username = match[1]
	} else {
		logger.Logf(logger.ERROR, "Failed to get username from log: %s", log)
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		logger.Logf(logger.ERROR, "Failed to set activity: %s for user %s", log, username)
		return
	}
	logger.InfoF("Successfully set activity: %s for user %s", log, username)
}
