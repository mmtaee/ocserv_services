package event

import (
	"context"
	"github.com/mmtaee/go-oc-utils/logger"
	"gorm.io/gorm"
	"sync"
)

type WorkerEvent struct {
	db        *gorm.DB
	eventChan chan *SchemaEvent
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	handler   RepositoryEventInterface
}

var Worker *WorkerEvent

// Set configs and create Worker
func Set(db *gorm.DB, bufferSize int) {
	c, cancel := context.WithCancel(context.Background())
	Worker = &WorkerEvent{
		eventChan: make(chan *SchemaEvent, bufferSize),
		ctx:       c,
		cancel:    cancel,
		handler:   NewEventRepository(db),
	}
}

// GetWorker return WorkerEvent
func GetWorker() *WorkerEvent {
	return Worker
}

// Start Workers
func (w *WorkerEvent) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.runWorker(i)
	}
	w.wg.Wait()
}

// Stop Workers
func (w *WorkerEvent) Stop() {
	logger.Log(logger.WARNING, "Stopping event workers")
	w.cancel()
	w.wg.Wait()
	close(w.eventChan)
	logger.Info("Event workers stopped")
}

// runWorker method run workers by workerID
func (w *WorkerEvent) runWorker(workerID int) {
	defer w.wg.Done()
	logger.InfoF("Event worker %d started", workerID+1)
	for {
		select {
		case e := <-w.eventChan:
			if err := w.handler.Apply(w.ctx, e); err != nil {
				logger.InfoF("Worker %d failed to process event: %v", workerID, err)
			} else {
				logger.InfoF("Worker %d processed event: %d", workerID, e.ID)
			}
		case <-w.ctx.Done():
			logger.Logf(logger.WARNING, "Worker %d shutting down...", workerID)
			return
		}
	}
}

// AddEvent method to start creating Event in channel
func (w *WorkerEvent) AddEvent(event *SchemaEvent) {
	select {
	case w.eventChan <- event:
		logger.InfoF("Added event: %d", event.ID)
	case <-w.ctx.Done():
		logger.InfoF("Event worker %s shutting down...", w.ctx.Err().Error())
	}
}
