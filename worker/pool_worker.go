package worker

import (
	"sync"
	"time"
)

// PoolWorker represents a worker in a pool
type PoolWorker struct {
	client  Client
	backoff Backoff
	channel chan int64
}

// NewPoolWorker creates a pool worker
func NewPoolWorker(client Client) *PoolWorker {
	return &PoolWorker{
		client:  client,
		channel: make(chan int64),
	}
}

// Run starts the pool worker
func (pw *PoolWorker) Run(handler ResponseHandler, wg *sync.WaitGroup) {
	for height := range pw.channel {
		pw.Process(height, handler)

		if wg != nil {
			wg.Done()
		}
	}
}

// Process handles the processing of a given height
func (pw *PoolWorker) Process(height int64, handler ResponseHandler) {
	err := pw.client.Send(Request{Height: height})
	if err != nil {
		pw.Reconnect()
		return
	}

	var res Response

	err = pw.client.Receive(&res)
	if err != nil {
		pw.Reconnect()
		return
	}

	handler(res)
}

// Reconnect reestablishes the connection with a worker
func (pw *PoolWorker) Reconnect() error {
	time.Sleep(pw.backoff.Delay())

	pw.backoff.Attempt()

	err := pw.client.Reconnect()
	if err != nil {
		return err
	}

	pw.backoff.Reset()

	return nil
}

// Stop stops the pool worker
func (pw *PoolWorker) Stop() {
	close(pw.channel)
}
