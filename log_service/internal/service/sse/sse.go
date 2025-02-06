package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Server is a basic struct for managing SSE connections.
type Server struct {
	Clients   map[chan string]string
	LogChan   chan string
	ClientMu  *sync.Mutex
	WebServer *http.Server
}

func NewSSEServer() *Server {
	return &Server{
		Clients:  make(map[chan string]string),
		LogChan:  make(chan string),
		ClientMu: &sync.Mutex{},
	}
}

func (server *Server) Start() {
	logger.Info("Log broadcast service started")
	go func() {
		for msg := range server.LogChan {
			server.Broadcast(msg)
		}
	}()

	go func() {
		server.StartWebService()
	}()
}

func (server *Server) StartWebService() {
	logger.Info("Log broadcast web service started")
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server.WebServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: http.HandlerFunc(server.ServerEventsHandler),
	}

	logger.InfoF("Starting server on %s:%s", host, port)
	if err := server.WebServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.CriticalF("ListenAndServe error: %v", err)
	}
}

// AddClient adds a new client connection.
func (server *Server) AddClient(client chan string, ip string) {
	server.ClientMu.Lock()
	defer server.ClientMu.Unlock()
	server.Clients[client] = ip
	logger.InfoF("Client %v (%s) connected", client, ip)
}

// RemoveClient removes a client connection.
func (server *Server) RemoveClient(client chan string) {
	server.ClientMu.Lock()
	defer server.ClientMu.Unlock()
	if ip, ok := server.Clients[client]; ok {
		logger.InfoF("Client %v (%s) disconnected", client, ip)
		delete(server.Clients, client)
		close(client)
	}
}

func (server *Server) Broadcast(msg string) {
	server.ClientMu.Lock()
	defer server.ClientMu.Unlock()
	for client := range server.Clients {
		client <- msg
	}
}

func checkToken(c context.Context, token string) error {
	var userToken models.UserToken
	db := database.Connection()
	err := db.WithContext(c).
		Table("user_tokens").
		Preload("User").Preload("User.Permission").
		Where("token = ? AND expire_at > ?", token, time.Now()).
		First(&userToken).Error
	if err != nil {
		return err
	}
	if userToken.User.IsAdmin {
		return nil
	}
	if !userToken.User.Permission.SeeServerLog {
		return errors.New("user not permitted")
	}
	return nil
}

func (server *Server) ServerEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	ctx := r.Context()
	queries := r.URL.Query()

	if len(queries["token"]) != 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid or missing 'token' query parameter",
			"code":  "Bad Request",
		})
		return
	}
	if err := checkToken(ctx, queries["token"][0]); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid 'token' or SeeServerLog permission not permitted",
			"code":  "Bad Request",
		})
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ch := make(chan string, 100)
	server.AddClient(ch, r.RemoteAddr)

	for {
		select {
		case message := <-ch:
			if message == "" {
				return
			}
			_, err := fmt.Fprintf(w, "data: %s\n\n", message)
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}
			flusher.Flush()
			time.Sleep(500 * time.Millisecond)
		case <-ctx.Done():
			return
		}
	}

}

// Shutdown closes all client channels.
func (server *Server) Shutdown() {
	server.ClientMu.Lock()
	defer server.ClientMu.Unlock()
	for ch, ip := range server.Clients {
		logger.InfoF("Closing client: %v (%s)", ch, ip)
		close(ch)
		delete(server.Clients, ch)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.WebServer.Shutdown(shutdownCtx); err != nil {
		logger.CriticalF("Server Shutdown error: %v", err)
	}
}
