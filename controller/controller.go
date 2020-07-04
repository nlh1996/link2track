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
	buffer []byte
	list   []string
	bs     []byte
	fspan  model.Span
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
	endCh chan bool
)

func init() {
	endCh = make(chan bool)
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
	getRes(env.URL)
	// fmt.Println("请求用时", time.Now().Sub(start))
}

func streamHandle() {
	size := env.StreamSize - 10000
	for {
		select {
		case <-endCh:
			fmt.Println("好难", len(model.ErrTid))
			for {
				span := <-model.Stream
				model.Mux.Lock()
				_, ok := model.ErrTid[span.Tid]
				model.Mux.Unlock()
				if ok {
					ws.WriteSpan(s2b(span.Data))
				}
				if len(model.Stream) == 0 {
					ws.WriteSpan([]byte("end"))
					return
				}
			}
		default:
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
	var res []byte
	for {
		n, err := resp.Body.Read(buffer)
		if n == 0 || err != nil {
			go filter(res)
			endCh <- true
			fmt.Println("读取结束", time.Now().Sub(start), n, err)
			//resp.Body.Close()
			return
		}
		res = append(res, buffer[:n]...)
		if len(res) > 100000000 {
			go filter(res)
			res = nil
		}
	}
	// if body, err := ioutil.ReadAll(resp.Body); err != nil {
	// 	log.Println(err)
	// } else {
	// 	go filter(body)
	// }
}

var count int

func filter(bs []byte) {
	st := time.Now()
	var res = false
	list = strings.Split(b2s(bs), sep)
	for _, v := range list {
		arr := strings.Split(v, sep2)
		if len(arr) < 9 {
			count++
			continue
		}
		fspan.Tid = arr[0]
		fspan.Data = v
		res = strings.Contains(arr[8], sep3)
		if res {
			ws.WriteTid(s2b(fspan.Tid))
			model.Mux.Lock()
			model.ErrTid[fspan.Tid] = ""
			model.Mux.Unlock()
			model.Stream <- fspan
			continue
		}
		res = strings.Contains(arr[8], sep4)
		if res {
			res = strings.Contains(arr[8], sep5)
			if !res {
				ws.WriteTid(s2b(fspan.Tid))
				model.Mux.Lock()
				model.ErrTid[fspan.Tid] = ""
				model.Mux.Unlock()
			}
		}
		model.Stream <- fspan
	}
	fmt.Println("计算用时", time.Now().Sub(st))
	fmt.Println("count=", count)
}
