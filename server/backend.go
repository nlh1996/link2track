package server

import (
	"bytes"
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

	"github.com/gin-gonic/gin"
)

var connPool map[string]net.Conn

// Server .
func Server() {
	connPool = make(map[string]net.Conn)
	// 创建监听
	listener, err := net.Listen("tcp", ":8003")
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
		connPool[addr] = conn
		//处理用户请求, 新建一个协程
		go readLoop(conn)
	}
}

func readLoop(conn net.Conn) {
	buf := make([]byte, 2000000) // 创建2048大小的缓冲区，用于read
	var (
		start  time.Time
		end    time.Time
		result string
		span   model.Span
		index  int
		arr    []string
	)
	for {
		//读取用户数据
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err = ", err)
			delete(connPool, conn.RemoteAddr().String())
			index = 0
			return
		}

		data := string(buf[:n])

		list := strings.Split(data, "\r")
		for _, v := range list {
			if v == "end" {
				index++
				if index == len(connPool) {
					start = time.Now()
					list2 := strings.Split(result, "\n")
					for _, item := range list2 {
						arr = strings.Split(item, "|")
						if len(arr) < 9 {
							continue
						}
						span.Tid = arr[0]
						span.Time = arr[1]
						span.Data = item + "\n"
						model.SpanMap[span.Tid] = append(model.SpanMap[span.Tid], span)
					}

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
					end = time.Now()
					fmt.Println("计算用时:", end.Sub(start))
					fmt.Println(model.Result)
					httpPost("http://localhost:" + env.ResPort + "/api/finished")
					return
				}
			}
			if len(v) < 20 {
				_, ok := model.ErrTid[v]
				if !ok {
					model.ErrTid[v] = ""
					v = v + "\r"
					for _, c := range connPool {
						c.Write([]byte(v))
					}
				}
			} else {
				result += v
			}
		}
	}
}

func postRes(url string) {
	requestBody := gin.H{"result": model.Result}
	data, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Invalid url for downloading: %s, error: %v", url, err)
	}
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := env.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
	if resp != nil {
		resp.Body.Close()
	}
}

func httpPost(URL string) {
	DataURLVal := url.Values{}
	for k, v := range model.Result {
		DataURLVal.Add(k, v)
	}
	data := "result=" + DataURLVal.Encode()
	resp, err := http.Post(URL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
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
