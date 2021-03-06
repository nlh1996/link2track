package controller

import (
	"cloud/env"
	"cloud/model"
	"cloud/utils"
	"cloud/ws"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	index  int
	start  time.Time
	sep    = "\n"
	sep2   = "|"
	sep3   = "error=1"
	sep4   = "code"
	sep5   = "code=200"
	// sep    = []byte("\n")
	// sep2   = []byte("|")
	// sep3   = []byte("error=1")
	// sep4   = []byte("code")
	// sep5   = []byte("code=200")
	b2s   = utils.Bytes2str
	s2b   = utils.Str2bytes
)

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
	} else {
		env.URL = "http://localhost:" + env.ResPort + "/api/finished"
	}
	c.String(200, "ok")
}

func startGet() {
	if env.Port == "8000" {
		env.URL = "http://localhost:" + env.ResPort + "/trace1.data"
	}
	if env.Port == "8001" {
		env.URL = "http://localhost:" + env.ResPort + "/trace2.data"
	}

	go streamHandle()

	start = time.Now()
	getTid()
	fmt.Println("第一次请求时间", time.Now().Sub(start))

	// start = time.Now()
	// getRes()
	// fmt.Println("第二次请求时间", time.Now().Sub(start))
}

func getTid() {
	req, err := http.NewRequest("GET", env.URL, nil)
	if err != nil {
		ws.WriteSpan(s2b("end"))
		log.Fatalf("Invalid url for downloading")
	}
	req.Header.Set("Accept-Charset", "utf-8")

	resp, err := env.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	readData(resp)
}

// getRes .
func getRes() {
	req, err := http.NewRequest("GET", env.URL, nil)
	if err != nil {
		ws.WriteSpan(s2b("end"))
		log.Fatalf("Invalid url for downloading: %s, error: %v", env.URL, err)
	}
	req.Header.Set("Accept-Charset", "utf-8")
	resp, err := env.Client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	readData2(resp)
}

func readData(resp *http.Response) {
	buffer := make([]byte, env.BufferSize)
	var res []byte
	for {
		n, err := resp.Body.Read(buffer)
		res = append(res, buffer[:n]...)
		if n == 0 || err != nil {
			go filter(res, 1)
			fmt.Println("读取结束", n, err)
			resp.Body.Close()
			return
		}
		if len(res) > 60000000 {
			go filter(res, 0)
			res = nil
		}
	}
}

func readData2(resp *http.Response) {
	buffer := make([]byte, env.BufferSize)
	var res []byte
	for {
		n, err := resp.Body.Read(buffer)
		res = append(res, buffer[:n]...)
		if n == 0 || err != nil {
			go filter2(res, 1)
			fmt.Println("读取结束", n, err)
			// resp.Body.Close()
			return
		}
		if len(res) > 60000000 {
			go filter2(res, 0)
			res = nil
		}
	}
}

func filter(bs []byte, i int) {
	st := time.Now()
	var res = false
	list := strings.Split(b2s(bs), sep)
	for _, v := range list {
		arr := strings.Split(v, sep2)
		if len(arr) < 9 || len(arr[0]) < 12 {
			continue
		}
		span := &model.Span{}
		span.Tid = arr[0]
		span.Data = v
		res = strings.Contains(arr[8], sep3)
		if res {
			ws.WriteTid(s2b(arr[0]))
			model.Stream <- span
			continue
		}
		res = strings.Contains(arr[8], sep4)
		if res {
			res = strings.Contains(arr[8], sep5)
			if !res {
				ws.WriteTid(s2b(arr[0]))
			}
		}
		model.Stream <- span
	}
	if i == 1 {
		model.EndSign = 1
	}
	fmt.Println("1计算用时", time.Now().Sub(st))
}

func filter2(bs []byte, i int) {
	st := time.Now()
	list := strings.Split(b2s(bs), sep)
	for _, v := range list {
		arr := strings.Split(v, sep2)
		if len(arr) < 9 || len(arr[0]) < 12 {
			fmt.Println(arr[0])
			continue
		}
		model.Mux.Lock()
		_, ok := model.ErrTid[arr[0]]
		model.Mux.Unlock()
		if ok {
			ws.WriteSpan(s2b(v))
		}
	}
	if i == 1 {
		ws.WriteSpan(s2b("end"))
	}
	fmt.Println("2计算用时", time.Now().Sub(st))
}

func streamHandle() {
	size := env.StreamSize - 1000
	for {
		if model.EndSign == 1 {
			for {
				span := <-model.Stream
				_, ok := model.ErrTid[span.Tid]
				if ok {
					ws.WriteSpan(s2b(span.Data))
				}
				if len(model.Stream) == 0 {
					ws.WriteSpan(s2b("end"))
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
				ws.WriteSpan(s2b(span.Data))
			}
		}
	}
}