package events

import (
	"fmt"
	"github.com/mmtaee/go-oc-utils/logger"
	"log"
	"net/http"
	"sync"
	"time"
)

func NewSSEServer() *SSEServer {
	return &SSEServer{
		Clients:  make(map[chan string]string),
		LogChan:  make(chan string),
		ClientMu: &sync.Mutex{},
	}
}

// SSEServer is a basic struct for managing SSE connections.
type SSEServer struct {
	Clients  map[chan string]string
	LogChan  chan string
	ClientMu *sync.Mutex
}

// AddClient adds a new client connection.
func (sse *SSEServer) AddClient(client chan string, ip string) {
	sse.ClientMu.Lock()
	defer sse.ClientMu.Unlock()
	sse.Clients[client] = ip
	logger.InfoF("Client %v (%s) connected", client, ip)
}

// RemoveClient removes a client connection.
func (sse *SSEServer) RemoveClient(client chan string) {
	sse.ClientMu.Lock()
	defer sse.ClientMu.Unlock()
	if ip, ok := sse.Clients[client]; ok {
		logger.InfoF("Client %v (%s) disconnected", client, ip)
		delete(sse.Clients, client)
		close(client)
	}
}

func (sse *SSEServer) Broadcast(msg string) {
	sse.ClientMu.Lock()
	defer sse.ClientMu.Unlock()
	for client := range sse.Clients {
		client <- msg
	}
}

func (sse *SSEServer) ServerEventsHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	if len(queries["token"]) != 1 {
		http.Error(w, "Invalid or missing 'token' query parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ch := make(chan string, 100)
	sse.AddClient(ch, r.RemoteAddr)
	defer sse.RemoveClient(ch)

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
			time.Sleep(1 * time.Second)
		case <-r.Context().Done():
			return
		}
	}
}

// CloseAllClients closes all client channels.
func (sse *SSEServer) CloseAllClients() {
	sse.ClientMu.Lock()
	defer sse.ClientMu.Unlock()
	for ch, ip := range sse.Clients {
		logger.InfoF("Closing client: %v (%s)", ch, ip)
		close(ch)
		delete(sse.Clients, ch)
	}
}
