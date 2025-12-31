package stream

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WsMessage struct {
	MessageType int
	Data        []byte
}

type WsConnection struct {
	wsSocket  *websocket.Conn
	inChan    chan *WsMessage
	outChan   chan *WsMessage
	closeChan chan struct{}
	closeOnce sync.Once
}

// WsConn 创建连接
func WsConn(response http.ResponseWriter, request *http.Request, header string) (*WsConnection, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
		Subprotocols: []string{
			header,
		},
	}

	wsSocket, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		return nil, err
	}

	conn := &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 256),
		outChan:   make(chan *WsMessage, 256),
		closeChan: make(chan struct{}),
	}

	// 启动读写 goroutine
	go conn.readLoop()
	go conn.writeLoop()

	return conn, nil
}

// ----------
// READ LOOP
// ----------
func (c *WsConnection) readLoop() {
	defer c.WsClose()

	for {
		msgType, data, err := c.wsSocket.ReadMessage()
		if err != nil {
			return
		}

		select {
		case c.inChan <- &WsMessage{msgType, data}:
		case <-c.closeChan:
			return
		}
	}
}

// ----------
// WRITE LOOP
// ----------
func (c *WsConnection) writeLoop() {
	defer c.WsClose()

	for {
		select {
		case msg := <-c.outChan:
			if err := c.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				return
			}

		case <-c.closeChan:
			return
		}
	}
}

// ----------
// Public API
// ----------

// WsWrite 非阻塞，防止 outChan 堵塞导致死锁
func (c *WsConnection) WsWrite(messageType int, data []byte) error {
	select {
	case c.outChan <- &WsMessage{messageType, data}:
		return nil
	case <-c.closeChan:
		return errors.New("websocket closed")
	default:
		// 避免写阻塞 — 不让 exec 卡死
		return errors.New("websocket write buffer full")
	}
}

func (c *WsConnection) WsRead() (*WsMessage, error) {
	select {
	case msg := <-c.inChan:
		return msg, nil
	case <-c.closeChan:
		return nil, errors.New("websocket closed")
	}
}

func (c *WsConnection) WsClose() {
	c.closeOnce.Do(func() {
		close(c.closeChan)
		c.wsSocket.Close()
	})
}
