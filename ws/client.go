package ws

import (
	"cloud/env"
	"cloud/model"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

var (
	ws1      *websocket.Conn
	ws2      *websocket.Conn
	err      error
	inChan1  = make(chan []byte, 10000000)
	outChan1 = make(chan []byte, 10000000)
	inChan2  = make(chan []byte, 10000000)
	outChan2 = make(chan []byte, 10000000)
)

// Dial .
func Dial() {
	go writeLoop()
	ws1 = nil
	ws2 = nil
	url1 := "ws://localhost:8002/?id=" + env.Port
	url2 := "ws://localhost:8002/?id=2"
	origin := "http://localhost:8002/"
	log.Println("客户端等待服务器连接中...")
	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)
		ws1, err = websocket.Dial(url1, "", origin)
		if ws1 != nil && err == nil {
			break
		}
	}
	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)
		ws2, err = websocket.Dial(url2, "", origin)
		if ws2 != nil && err == nil {
			break
		}
	}

	log.Println("成功连接...")
	go read()
}

// WriteTid .
func WriteTid(data []byte) {
	outChan1 <- data
}

// WriteSpan .
func WriteSpan(data []byte) {
	outChan2 <- data
}

func read() {
	var data = make([]byte, 1024)
	for {
		m, err := ws1.Read(data)
		if err != nil {
			log.Println(err)
			return
		}
		key := string(data[:m])
		model.Mux.Lock()
		model.ErrTid[key] = ""
		model.Mux.Unlock()
	}
}

func writeLoop() {
	for {
		select {
		case data := <-outChan1:
			if ws1 != nil {
				_, err = ws1.Write(data)
				if err != nil {
					log.Println(err)
				}
			}
		case data := <-outChan2:
			if ws2 == nil {
				_, err = ws2.Write(data)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
