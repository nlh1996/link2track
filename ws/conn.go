package ws

import (
	"cloud/model"
	"context"
	"errors"
	"fmt"

	"log"

	"github.com/gorilla/websocket"
)

// Connection .
type Connection struct {
	ID        string
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	ctx       context.Context
	cancel    context.CancelFunc
	connected bool
}

// NewConnection .
func NewConnection(wsConn *websocket.Conn, id string) (*Connection, error) {
	conn := &Connection{
		wsConnect: wsConn,
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		ID:        id,
		connected: true,
	}
	conn.ctx, conn.cancel = context.WithCancel(context.Background())
	if id != "2" {
		p := GetConnPool()
		p.Set(conn)
	}
	return conn, nil
}

// Start .
func (conn *Connection) Start() (data []byte, err error) {
	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()

	for {
		select {
		case data = <-conn.inChan:
			if conn.ID != "2" {
				fmt.Println(string(data))
			}

		case <-conn.ctx.Done():
			return
		}
	}
}

// Close .
func (conn *Connection) Close() {
	if conn.connected {
		conn.wsConnect.Close()
		conn.cancel()
		conn.connected = false
		GetConnPool().DelByID(conn.ID)
		log.Println("连接", conn.ID, "已经关闭！！！")
	}
}

// Send .
func (conn *Connection) Send(msgBytes []byte) (err error) {
	fmt.Println("SEND->", string(msgBytes))

	select {
	case conn.outChan <- msgBytes:
	case <-conn.ctx.Done():
		err = errors.New("connection is closed")
	}
	return
}

// 内部实现
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		d := string(data)
		if conn.ID != "2" {
			model.Mux.Lock()
			_, ok := model.ErrTid[d]
			model.Mux.Unlock()
			if !ok {
				model.Mux.Lock()
				model.ErrTid[d] = ""
				model.Mux.Unlock()
				for _, c := range GetConnPool().Pool {
					c.Send(data)
				}
			}
		}
	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
			if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
				goto ERR
			}
		case <-conn.ctx.Done():
			goto ERR
		}
	}

ERR:
	conn.Close()
}
