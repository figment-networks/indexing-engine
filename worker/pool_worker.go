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
		pw.process(height, handler)

		if wg != nil {
			wg.Done()
		}
	}
}

// process handles the processing of a given height
func (pw *PoolWorker) process(height int64, handler ResponseHandler) {
	err := pw.client.Send(Request{Height: height})
	if err != nil {
		pw.reconnect()
		return
	}

	var res Response

	err = pw.client.Receive(&res)
	if err != nil {
		pw.reconnect()
		return
	}

	handler(res)
}

// reconnect reestablishes the connection with a worker
func (pw *PoolWorker) reconnect() error {
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
