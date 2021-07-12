package worker

// Request respresents a manager request
type Request struct {
	Height int64
}

// Response represents a worker response
type Response struct {
	Height  int64
	Success bool
	Error   string
}

// Client interacts with a worker
type Client interface {
	Send(Request) error
	Receive(*Response) error
	Close() error
	Reconnect() error
}

// Server interacts with a manager
type Server interface {
	Receive(*Request) error
	Send(Response) error
}

// RequestHandler handles a request from a manager
type RequestHandler func(Request) error

// ResponseHandler handles a response from a worker
type ResponseHandler func(Response)
