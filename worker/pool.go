package worker

import (
	"reflect"
	"sync"
)

// Pool represents a pool of workers
type Pool struct {
	workers []*PoolWorker
	wg      sync.WaitGroup
}

// AddWorker adds a worker to the pool
func (p *Pool) AddWorker(worker *PoolWorker) {
	p.workers = append(p.workers, worker)
}

// Run starts the worker pool
func (p *Pool) Run(handler ResponseHandler) {
	for _, worker := range p.workers {
		go worker.Run(handler, &p.wg)
	}
}

// Process schedules the processing of a given height
func (p *Pool) Process(height int64) {
	p.wg.Add(1)

	cases := make([]reflect.SelectCase, len(p.workers))

	for i, worker := range p.workers {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(worker.channel),
			Send: reflect.ValueOf(height),
		}
	}

	reflect.Select(cases)
}

// Wait blocks until all workers are finished
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Stop stops the worker pool
func (p *Pool) Stop() {
	for _, worker := range p.workers {
		worker.Stop()
	}
}
