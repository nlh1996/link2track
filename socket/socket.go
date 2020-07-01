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
	conn net.Conn
	mux  sync.Mutex
)

// Init .
func Init() {
	// 主动发起连接请求
	conn = nil
	log.Println("客户端等待服务器连接中...")
	for conn == nil {
		conn, _ = net.Dial("tcp", "127.0.0.1:8003")
	}
	log.Println("成功连接...")
	// defer conn.Close() // 结束时，关闭连接
	go readLoop()
}

func readLoop() {
	// 接收服务端数据
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
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

// Write .
func Write(data string) {
	if conn == nil {
		return
	}
	data = data + "\r"
	// 发送数据
	_, err := conn.Write([]byte(data))
	if err != nil {
		log.Println("Write err:", err)
	}
}
