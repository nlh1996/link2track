package ws

import (
	"context"
	"errors"
	"fmt"
	"pojing/model"

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
	router    *Router
	connected bool
}

var (
	users *model.UserMap
)

func init() {
	users = model.GetUserMgr()
}

// NewConnection .
func NewConnection(wsConn *websocket.Conn, r *Router, id string) (*Connection, error) {
	conn := &Connection{
		wsConnect: wsConn,
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		router:    r,
		ID:        id,
		connected: true,
	}
	conn.ctx, conn.cancel = context.WithCancel(context.Background())
	p := GetConnPool()
	p.Set(conn)

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
			req := &Request{
				conn: conn,
				ID:  string(data[:4]),
				ByteData: data[4:],
			}
			// 请求跟路由绑定，并发处理请求
			go conn.router.DoMsgHandle(req)

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
		//阻塞在这里，等待inChan有空闲位置
		select {
		case conn.inChan <- data:
		case <-conn.ctx.Done(): // closeChan 感知 conn断开
			goto ERR
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
