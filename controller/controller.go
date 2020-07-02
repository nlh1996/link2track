package controller

import (
	"cloud/env"
	"cloud/model"
	"cloud/socket"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	buffer []byte
	list   []string
	bytes  []byte
	length int
	temp   string
	index  int
	start  time.Time
	end    time.Time
	fspan  model.Span
)

func init() {
	buffer = make([]byte, env.BufferSize)
}

// Ready .
func Ready(c *gin.Context) {
	c.String(200, "ok")
}

// SetParameter .
func SetParameter(c *gin.Context) {
	index++
	env.ResPort = c.Query("port")
	if env.Port != "8002" && index == 1 {
		go startGet()
	}
	c.String(200, "ok")
}

func startGet() {
	// if env.Port == "8000" {
	// 	env.URL = "http://localhost:" + env.ResPort + "/trace1.data"
	// }
	// if env.Port == "8001" {
	// 	env.URL = "http://localhost:" + env.ResPort + "/trace2.data"
	// }

	go byteStreamHandle()
	go streamHandle()

	start = time.Now()
	getRes(env.URL)
}

func byteStreamHandle() {
	for {
		select {
		case bytes = <-model.ByteStream:
			if string(bytes) == "end" {
				model.EndSign = 1
				fmt.Println(count)
				return
			}
			list = strings.Split(string(bytes), "\n")
			length = len(list)
			temp = list[length-1]
			list[0] = temp + list[0]
			filter(list[:length-1])
		}
	}
}

// 1514698
// 1514908
// 1512069

func streamHandle() {
	size := env.StreamSize - 1000
	for {
		if model.EndSign == 1 {
			for {
				span := <-model.Stream
				model.Mux.Lock()
				_, ok := model.ErrTid[span.Tid]
				model.Mux.Unlock()
				if ok {
					socket.Write(span.Data)
				}
				if len(model.Stream) == 0 {
					socket.Write("end")
					return
				}
			}
		}
		if len(model.Stream) > size {
			span := <-model.Stream
			model.Mux.Lock()
			_, ok := model.ErrTid[span.Tid]
			model.Mux.Unlock()
			if ok {
				socket.Write(span.Data)
			}
		}
	}
}

// getRes .
func getRes(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Invalid url for downloading: %s, error: %v", url, err)
	}
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := env.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	readData(resp)
}

func readData(resp *http.Response) {
	for {
		n, err := resp.Body.Read(buffer)
		if n == 0 || err != nil {
			model.ByteStream <- []byte("end")
			end = time.Now()
			fmt.Println("读取结束", end.Sub(start), n, err)
			//resp.Body.Close()
			return
		}
		// fmt.Println(n)
		model.ByteStream <- buffer[:n]
	}
}

var count int

func filter(list []string) {
	var res = false
	for _, v := range list {
		arr := strings.Split(v, "|")
		fspan.Tid = arr[0]
		fspan.Data = v + "\n"
		model.Stream <- fspan
		l := len(arr) - 1
		res = strings.Contains(arr[l], "error=1")
		if res {
			socket.Write(arr[0])
			continue
		}
		res = strings.Contains(arr[l], "code")
		if res {
			res = strings.Contains(arr[l], "code=200")
			if !res {
				socket.Write(arr[0])
			}
		}
	}
}
