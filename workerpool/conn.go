package workerpool

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MessageHandler func([]byte) error
type ErrorHandler func(error)

type connection struct {
	id     string
	socket *websocket.Conn
	mutex  *sync.Mutex
	closed bool

	onMessageHandler MessageHandler
	onError          ErrorHandler
}

type Connection interface {
	ID() string

	OnMessage(handler MessageHandler)
	OnError(handler ErrorHandler)

	Hold(ctx context.Context)
	Send(message string) error
}

func (c *connection) ID() string {
	return c.id
}

func (c *connection) OnMessage(handler MessageHandler) {
	c.onMessageHandler = handler
}

func (c *connection) OnError(handler ErrorHandler) {
	c.onError = handler
}

func (c *connection) Hold(ctx context.Context) {

}

func (c *connection) Send(message string) error {
	return nil
}

func (c *connection) close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.closed {
		fmt.Println("close websocket connection")
		if err := c.socket.Close(); err != nil {
			fmt.Println("error close ws", err)
		}
		c.closed = true
	}
}

func (c *connection) addSetting() {
	c.socket.SetReadLimit(int64(maxMessageSize))
	c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	c.socket.SetCloseHandler(func(code int, msg string) error {
		fmt.Println("socket close handler", code, msg)
		message := websocket.FormatCloseMessage(code, "")
		c.socket.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		c.close()
		return nil
	})
}

func NewSocketConnection(writer http.ResponseWriter, req *http.Request) (Connection, error) {
	var (
		upgrader = websocket.Upgrader{
			HandshakeTimeout: 60 * time.Second,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: true,
		}

		responseHeader http.Header
	)

	socket, err := upgrader.Upgrade(writer, req, responseHeader)
	if err != nil {
		return nil, fmt.Errorf("cannot upgrade websocket connection %w", err)
	}

	conn := &connection{
		socket: socket,
		mutex:  new(sync.Mutex),
		closed: false,
	}

	conn.addSetting()

	return conn, nil
}
