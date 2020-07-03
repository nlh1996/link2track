package server

import (
	"cloud/env"
	"cloud/model"
	"cloud/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var (
	connPool map[string]net.Conn
	ch       chan string
	res      string
)

// Server .
func Server(port string) {
	// 创建监听
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("listen err:", err)
		return
	}
	defer listener.Close() // 主协程结束时，关闭listener

	log.Println("服务器等待客户端建立连接...")
	// 等待客户端连接请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept err = ", err)
			return
		}
		addr := conn.RemoteAddr().String()
		log.Println(addr, "connected.")
		if port == ":8003" {
			ch = make(chan string)
			go read()
			//处理用户请求, 新建一个协程
			go readLoop(conn)
		} else {
			connPool = make(map[string]net.Conn)
			connPool[addr] = conn
			go readTid(conn)
		}
	}
}

func readTid(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		//读取用户数据
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err = ", err)
			delete(connPool, conn.RemoteAddr().String())
			return
		}

		list := strings.Split(string(buf[:n]), "\r")
		for _, v := range list {
			if len(v) < 20 && v != "" {
				model.Mux.Lock()
				_, ok := model.ErrTid[v]
				model.Mux.Unlock()
				if !ok {
					model.Mux.Lock()
					model.ErrTid[v] = ""
					model.Mux.Unlock()
					v = v + "\r"
					for _, c := range connPool {
						c.Write([]byte(v))
					}
				}
			}
		}
	}
}

func read() {
	var i int
	for {
		select {
		case str := <-ch:
			i++
			res += str
			if i == len(connPool) {
				handle()
				break
			}
		}
	}
}

func handle() {
	start := time.Now()
	var (
		span model.Span
		arr  []string
	)
	list := strings.Split(res, "\r")
	for _, item := range list {
		arr = strings.Split(item, "|")
		if len(arr) < 2 {
			fmt.Println(item)
			continue
		}
		if arr[0] == "" {
			fmt.Println(item)
		}
		span.Tid = arr[0]
		span.Time = arr[1]
		span.Data = item
		model.SpanMap[span.Tid] = append(model.SpanMap[span.Tid], span)
	}

	// 排序
	for k, s := range model.SpanMap {
		// if k == "" {
		// 	fmt.Println(k, s)
		// }
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
	fmt.Println(len(model.Result))
	fmt.Println(len(model.ErrTid))
	httpPost("http://localhost:" + env.ResPort + "/api/finished")
}

func readLoop(conn net.Conn) {
	buf := make([]byte, 4096) // 创建2048大小的缓冲区，用于read
	var (
		result string
		//list   []string
		index  int
	)
	for {
		index ++
		//读取用户数据
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}

		if string(buf[:n]) == "end\r" {
			ch <- result
			break
		}
		
		result += string(buf[:n])
		if index == 100 {
			fmt.Println(result)
		}
		// list = strings.Split(string(buf[:n]), "\r")
		// if index < 40 {
		// 	fmt.Println(string(buf[:n]))
		// }
		// for _, v := range list {
		// 	if v == "end" {
		// 		ch <- result
		// 		break
		// 	}
		// 	result += v
		// }
	}
}

func httpPost(URL string) {
	DataURLVal := url.Values{}
	mjson, _ := json.Marshal(model.Result)
	mString := string(mjson)
	DataURLVal.Add("result", mString)
	resp, err := http.Post(URL,
		"application/x-www-form-urlencoded",
		strings.NewReader(DataURLVal.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(body))
	if resp != nil {
		resp.Body.Close()
	}
}
