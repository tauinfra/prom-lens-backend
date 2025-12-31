package ws

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MaxQueueSize = 1000
)

var ErrConnClosed = errors.New("websocket connection closed")

type Connection struct {
	ws   *websocket.Conn
	in   chan []byte
	out  chan []byte
	done chan struct{}

	once sync.Once
}

// NewWS creates a websocket connection wrapper
func NewWS(ws *websocket.Conn) *Connection {
	c := &Connection{
		ws:   ws,
		in:   make(chan []byte, MaxQueueSize),
		out:  make(chan []byte, MaxQueueSize),
		done: make(chan struct{}),
	}

	go c.readLoop()
	go c.writeLoop()

	return c
}

// Close closes the connection safely (idempotent)
func (c *Connection) Close() {
	c.once.Do(func() {
		close(c.done)
		_ = c.ws.Close()
	})
}

// Read reads one message from websocket
func (c *Connection) Read() ([]byte, error) {
	select {
	case data := <-c.in:
		return data, nil
	case <-c.done:
		return nil, ErrConnClosed
	}
}

// Write sends one message to websocket
func (c *Connection) Write(data []byte) error {
	select {
	case c.out <- data:
		return nil
	case <-c.done:
		return ErrConnClosed
	}
}

// readLoop reads from websocket → in channel
func (c *Connection) readLoop() {
	defer c.Close()

	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			return
		}

		select {
		case c.in <- data:
		case <-c.done:
			return
		}
	}
}

// writeLoop writes from out channel → websocket
func (c *Connection) writeLoop() {
	defer c.Close()

	for {
		select {
		case data := <-c.out:
			if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-c.done:
			return
		}
	}
}
