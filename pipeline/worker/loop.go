package worker

import (
	"errors"
	"io"
)

// Loop represents the processing loop of a worker
type Loop struct {
	server Server
}

// NewLoop creates a worker loop
func NewLoop(server Server) *Loop {
	return &Loop{server: server}
}

// Run starts the worker loop
func (l *Loop) Run(handler RequestHandler) {
	for {
		var req Request

		err := l.server.Receive(&req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			panic(err)
		}

		err = handler(req)

		var msg string
		if err != nil {
			msg = err.Error()
		}

		l.server.Send(Response{
			Height:  req.Height,
			Success: err == nil,
			Error:   msg,
		})
	}
}
