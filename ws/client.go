package ws

import (
	"fmt"
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
	url1 := "http://localhost:8002/"
	//url2 := "http://localhost:8002/?id=2"
	origin := "ws://localhost:8002/"
	log.Println("客户端等待服务器连接中...")
	// for ws1 == nil {
	ws, err := websocket.Dial(url1, "", origin)
	if err != nil {
		log.Panicln(err)
	}
	ws.Write([]byte("111"))
	// }
	// for ws2 == nil {
	// 	ws2, _ = websocket.Dial(url2, "", origin)
	// }
	log.Println("成功连接...")
	go read()
}

// WriteTid .
func WriteTid(data string) {
	_, err = ws1.Write([]byte(data))
	if err != nil {
		log.Println(err)
	}
}

// WriteSpan .
func WriteSpan(data string) {
	_, err = ws2.Write([]byte(data))
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
		}
		fmt.Printf("Receive: %s\n", data[:m])
	}
}
