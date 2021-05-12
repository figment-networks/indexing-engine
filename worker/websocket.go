package worker

import (
	"strings"

	"golang.org/x/net/websocket"
)

// WebsocketClient interacts with a worker using a websocket
type WebsocketClient struct {
	url  string
	conn *websocket.Conn
}

var _ Client = (*WebsocketClient)(nil)

// NewWebsocketClient creates a websocket client
func NewWebsocketClient(endpoint string) (*WebsocketClient, error) {
	url := endpoint
	if !strings.HasPrefix(url, "ws://") {
		url = "ws://" + url
	}

	wc := WebsocketClient{url: url}

	if err := wc.Connect(); err != nil {
		return nil, err
	}

	return &wc, nil
}

// Connect establishes a websocket connection
func (wc *WebsocketClient) Connect() error {
	conn, err := websocket.Dial(wc.url, "", "http://localhost")
	if err != nil {
		return err
	}

	wc.conn = conn

	return nil
}

// Send sends a request to a worker
func (wc *WebsocketClient) Send(req Request) error {
	return websocket.JSON.Send(wc.conn, req)
}

// Receive receives a response from a worker
func (wc *WebsocketClient) Receive(res *Response) error {
	return websocket.JSON.Receive(wc.conn, res)
}

// Close closes the websocket connection
func (wc *WebsocketClient) Close() error {
	return wc.conn.Close()
}

// Reconnect reestablishes a websocket connection
func (wc *WebsocketClient) Reconnect() error {
	wc.Close()
	return wc.Connect()
}

// WebsocketServer interacts with a manager using a websocket
type WebsocketServer struct {
	conn *websocket.Conn
}

var _ Server = (*WebsocketServer)(nil)

// NewWebsocketServer creates a websocket server
func NewWebsocketServer(conn *websocket.Conn) *WebsocketServer {
	return &WebsocketServer{conn: conn}
}

// Receive receives a request from a manager
func (ws *WebsocketServer) Receive(req *Request) error {
	return websocket.JSON.Receive(ws.conn, req)
}

// Send sens a response back to a manager
func (ws *WebsocketServer) Send(res Response) error {
	return websocket.JSON.Send(ws.conn, res)
}

// Close closes the websocket connection
func (ws *WebsocketServer) Close() error {
	return ws.conn.Close()
}
