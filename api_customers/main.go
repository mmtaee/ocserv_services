package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"

	"github.com/gorilla/handlers"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type RequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Response struct {
	User              models.OcUser                    `json:"oc_user"`
	Activities        []models.OcUserActivity          `json:"activities"`
	TrafficStatistics []models.OcUserTrafficStatistics `json:"traffic_statistics"`
}

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.Connect(&database.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
	}, false)

	httpRateLimiter := rateLimiter()

	mux := http.NewServeMux()
	mux.Handle("/client/dashboard/", handlers.LoggingHandler(os.Stdout, httpRateLimiter.RateLimit(http.HandlerFunc(Dashboard))))
	server := &http.Server{Addr: ":8080", Handler: handlers.CompressHandler(mux)}
	go func() {
		logger.InfoF("Starting server on %s:%s", host, port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logf(logger.ERROR, "Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	fmt.Println()
	logger.InfoF("Signal received: %s. Initiating shutdown...", sig)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.CriticalF(" Server forced to shutdown: %v\n", err)
	}
	database.Close()
	logger.Info("Database close successfully")
	logger.Info("Server shutdown complete")
}

func rateLimiter() throttled.HTTPRateLimiterCtx {
	store, err := memstore.NewCtx(65536)
	if err != nil {
		logger.CriticalF("Error creating memstore: %v", err)
	}
	quota := throttled.RateQuota{
		MaxRate:  throttled.PerHour(10),
		MaxBurst: 1,
	}
	rater, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		logger.CriticalF("Error creating rate limiter: %v", err)
	}
	return throttled.HTTPRateLimiterCtx{
		RateLimiter: rater,
		VaryBy:      &throttled.VaryBy{Path: true},
		DeniedHandler: http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Max request limit exceeded", 429)
		})),
	}
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.CriticalF("Failed to close request body: %s", err)
		}
	}(r.Body)

	var requestBody RequestBody
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Failed to parse request body. Username and password not found", http.StatusBadRequest)
		return
	}

	db := database.Connection()

	var ocUser models.OcUser
	err = db.WithContext(r.Context()).
		Where("username = ? AND password = ?", requestBody.Username, requestBody.Password).
		Find(&ocUser).Error
	if err != nil {
		http.Error(w, fmt.Sprintf("User with username(%s) not found", requestBody.Username), http.StatusNotFound)
		return
	}

	var activities []models.OcUserActivity
	err = db.WithContext(r.Context()).Where("oc_user_id = ?", ocUser.ID).Find(&activities).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("failed to get user activities. error: %v", err), http.StatusNotFound)
			return
		}
	}

	var trafficStatistics []models.OcUserTrafficStatistics
	err = db.WithContext(r.Context()).Where("oc_user_id = ?", ocUser.ID).Find(&trafficStatistics).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("failed to get user traffic statistics. error: %v", err), http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{
		User:              ocUser,
		Activities:        activities,
		TrafficStatistics: trafficStatistics,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
