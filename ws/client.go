package ws

import (
	"cloud/env"
	"cloud/model"
	"log"

	"golang.org/x/net/websocket"
)

var (
	ws1 *websocket.Conn
	ws2 *websocket.Conn
	err error
)

// Dial .
func Dial() {
	url1 := "ws://localhost:8002/?id=" + env.Port
	url2 := "ws://localhost:8002/?id=2"
	origin := "http://localhost:8002/"
	log.Println("客户端等待服务器连接中...")
	for ws1 == nil {
		ws1, _ = websocket.Dial(url1, "", origin)
	}
	for ws2 == nil {
		ws2, _ = websocket.Dial(url2, "", origin)
	}
	log.Println("成功连接...")
	go read()
}

// WriteTid .
func WriteTid(data []byte) {
	_, err = ws1.Write(data)
	if err != nil {
		log.Println(err)
	}
}

// WriteSpan .
func WriteSpan(data []byte) {
	_, err = ws2.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func read() {
	var data = make([]byte, 512)
	for {
		m, err := ws1.Read(data)
		if err != nil {
			log.Println(err)
			return
		}
		// fmt.Print(" ", utils.Bytes2str(data[:m]))
		key := string(data[:m])
		model.ErrTid[key] = ""
	}
}
