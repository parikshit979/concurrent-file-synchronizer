package tasks

import (
	"context"
	"sync"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/tasks/task"
)

type WorkerPool struct {
	workerCount int
	taskQueue   chan task.Task
	wg          *sync.WaitGroup
	errors      chan error
	ctx         context.Context
	cancel      context.CancelFunc
	eventBus    *eventbus.EventBus
}

func NewWorkerPool(workerCount int, eventBus *eventbus.EventBus) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan task.Task),
		wg:          &sync.WaitGroup{},
		errors:      make(chan error, 100),
		ctx:         ctx,
		cancel:      cancel,
		eventBus:    eventBus,
	}
}

func (w *WorkerPool) Submit(task task.Task) {
	w.wg.Add(1)
	w.taskQueue <- task
}

func (w *WorkerPool) worker() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case task, ok := <-w.taskQueue:
			if !ok {
				return
			}
			if err := task.Execute(w.ctx, w.eventBus); err != nil {
				w.errors <- err
			}
			w.wg.Done()
		}
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.workerCount; i++ {
		go w.worker()
	}
}

func (w *WorkerPool) Wait() []error {
	w.wg.Wait()
	w.cancel()
	close(w.errors)
	var errs []error
	for err := range w.errors {
		errs = append(errs, err)
	}
	return errs
}

func (w *WorkerPool) Stop() {
	w.cancel()
	close(w.taskQueue)
}
