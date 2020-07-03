package socket

import (
	"cloud/model"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	conn1 net.Conn
	conn2 net.Conn
	mux  sync.Mutex
)

// Init .
func Init() {
	// 主动发起连接请求
	conn1 = nil
	conn2 = nil
	log.Println("客户端等待服务器连接中...")
	for conn1 == nil {
		conn1, _ = net.Dial("tcp", "127.0.0.1:8003")
	}
	for conn2 == nil {
		conn2, _ = net.Dial("tcp", "127.0.0.1:8004")
	}
	log.Println("成功连接...")
	// defer conn1.Close() // 结束时，关闭连接
	go readLoop()
}

func readLoop() {
	// 接收服务端数据
	buf := make([]byte, 1024)
	for {
		n, err := conn2.Read(buf)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}

		list := strings.Split(string(buf[:n]), "\r")
		for _, v := range list {
			model.Mux.Lock()
			model.ErrTid[v] = ""
			model.Mux.Unlock()
		}
	}
}

// Write1 .
func Write1(data string) {
	if conn1 == nil {
		return
	}
	data = data + "\r"

	// 发送数据
	_, err := conn1.Write([]byte(data))
	if err != nil {
		log.Println("Write err:", err)
	}
}

// Write2 .
func Write2(data string) {
	if conn2 == nil {
		return
	}
	data = data + "\r"
	// 发送数据
	_, err := conn2.Write([]byte(data))
	if err != nil {
		log.Println("Write err:", err)
	}
}
