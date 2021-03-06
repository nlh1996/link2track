package ws

import (
	"bytes"
	"cloud/model"
	"cloud/utils"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
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
		inChan:    make(chan []byte, 10000000),
		outChan:   make(chan []byte, 1000000),
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

	var span model.Span

	for {
		select {
		case data := <-conn.inChan:
			d := utils.Bytes2str(data)
			if d == "end" {
				ch <- 0
			}
			if conn.ID != "2" && d != "" {
				model.Mux.Lock()
				_, ok := model.ErrTid[d]
				model.Mux.Unlock()
				if !ok {
					model.Mux.Lock()
					model.ErrTid[d] = ""
					model.Mux.Unlock()
					for _, c := range GetConnPool().Pool {
						mx.Lock()
						c.Send(utils.Str2bytes(d))
						mx.Unlock()
					}
				}
			} else {
				arr := strings.Split(d, "|")
				if len(arr) < 9 || arr[0] == "" {
					fmt.Println(d)
					continue
				}
				span.Tid = arr[0]
				span.Time = arr[1]
				span.Data = d + "\n"
				model.Mux.Lock()
				model.SpanMap[span.Tid] = append(model.SpanMap[span.Tid], span)
				model.Mux.Unlock()
			}

		case <-conn.ctx.Done():
			return
		}
	}
}

// Close .
func (conn *Connection) Close() {
	if conn.connected {
		ch <- 0
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
var mx sync.Mutex

func init() {
	mx = sync.Mutex{}
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

func handle() {
	start := time.Now()
	// 排序
	for k, s := range model.SpanMap {
		sort.Sort(s)
		model.Mux.Lock()
		model.SpanMap[k] = s
		model.Mux.Unlock()
	}
	// 聚合
	for k, s := range model.SpanMap {
		var buffer bytes.Buffer
		for _, item := range s {
			buffer.WriteString(item.Data)
		}
		// md5加密
		model.Result[k] = utils.Md5(buffer.String())
	}
	end := time.Now()
	log.Println("计算用时：", end.Sub(start))
	utils.HTTPPost()
	fmt.Println(len(model.Result))
	fmt.Println(len(model.ErrTid))
}
