package ws

import (
	"cloud/model"
	"cloud/utils"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

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
func (conn *Connection) Start() {
	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()
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
	select {
	case conn.outChan <- msgBytes:
	case <-conn.ctx.Done():
		err = errors.New("connection is closed")
	}
	return
}

var ch chan int

func init() {
	ch = make(chan int)
	go do()
}

func do() {
	var i int
	for {
		select {
		case <-ch:
			i++
			if i == len(GetConnPool().Pool) {
				handle()
				break
			}
		}
	}
}

// 内部实现
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
		arr  []string
		span model.Span
	)
	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		d := utils.Bytes2str(data)
		if d == "end" {
			ch <- 0
		}
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
		} else {
			arr = strings.Split(d, "|")
			if len(arr) < 2 {
				fmt.Println(d)
				continue
			}
			if arr[0] == "" {
				fmt.Println(d)
			}
			span.Tid = arr[0]
			span.Time = arr[1]
			span.Data = d + "\n"
			model.Mux.Lock()
			model.SpanMap[span.Tid] = append(model.SpanMap[span.Tid], span)
			model.Mux.Unlock()
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

func handle() {
	start := time.Now()
	// 排序
	for k, s := range model.SpanMap {
		sort.Sort(s)
		model.SpanMap[k] = s
	}
	// 聚合
	for k, s := range model.SpanMap {
		var str string
		for _, item := range s {
			str = str + item.Data
		}
		// md5加密
		model.Result[k] = utils.Md5(str)
	}
	end := time.Now()
	log.Println("计算用时：", end.Sub(start))
	fmt.Println(model.Result)
	fmt.Println(len(model.ErrTid))
	utils.HTTPPost()
}
