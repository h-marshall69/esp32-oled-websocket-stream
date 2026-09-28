package websocket

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	gorilla "github.com/gorilla/websocket"

	"oled/internal/domain"
)

const (
	writeTimeout = 3 * time.Second
	pongWait     = 60 * time.Second
	pingPeriod   = 45 * time.Second
)

var clientSequence atomic.Uint64

type Client struct {
	id   string
	conn *gorilla.Conn

	writeMu   sync.Mutex
	closeOnce sync.Once
}

func NewClient(conn *gorilla.Conn) *Client {
	client := &Client{
		id:   fmt.Sprintf("ws-%d", clientSequence.Add(1)),
		conn: conn,
	}

	conn.SetReadLimit(domain.MaxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	return client
}

func (c *Client) ID() string {
	if c == nil {
		return ""
	}
	return c.id
}

func (c *Client) ReadMessage() (int, []byte, error) {
	if c == nil || c.conn == nil {
		return 0, nil, errors.New("websocket not connected")
	}

	return c.conn.ReadMessage()
}

func (c *Client) SendText(message string) error {
	return c.writeMessage(gorilla.TextMessage, []byte(message))
}

func (c *Client) SendBinary(data []byte) error {
	return c.writeMessage(gorilla.BinaryMessage, data)
}

func (c *Client) writeMessage(messageType int, data []byte) error {
	if c == nil || c.conn == nil {
		return errors.New("websocket not connected")
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	_ = c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	return c.conn.WriteMessage(messageType, data)
}

func (c *Client) StartHeartbeat(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.sendPing(); err != nil {
					_ = c.Close()
					return
				}
			}
		}
	}()
}

func (c *Client) sendPing() error {
	if c == nil || c.conn == nil {
		return errors.New("websocket not connected")
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteControl(
		gorilla.PingMessage,
		nil,
		time.Now().Add(writeTimeout),
	)
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}

	var err error
	c.closeOnce.Do(func() {
		err = c.conn.Close()
	})

	return err
}
