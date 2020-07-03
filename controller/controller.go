package controller

import (
	"bytes"
	"cloud/env"
	"cloud/model"
	"cloud/utils"
	"cloud/ws"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	buffer []byte
	list   [][]byte
	bs     []byte
	start  time.Time
	end    time.Time
	fspan  model.Span
	index  int
	sep    = []byte("\n")
	sep2   = []byte("|")
	sep3   = []byte("error=1")
	sep4   = []byte("code")
	sep5   = []byte("code=200")
	b2s    = utils.Bytes2str
	s2b    = utils.Str2bytes
	endCh  chan bool
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

	go byteStreamHandle()
	go streamHandle()

	start = time.Now()
	getRes(env.URL)
}

func streamHandle() {
	size := env.StreamSize - 1000
	for {
		select {
		case <-endCh:
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
	for {
		n, err := resp.Body.Read(buffer)
		if n == 0 || err != nil {
			endCh <- true
			end = time.Now()
			fmt.Println("读取结束", end.Sub(start), n, err)
			//resp.Body.Close()
			return
		}
		model.ByteStream <- buffer[:n]
		// fmt.Println(b2s(buffer[:n]))
	}
}

func byteStreamHandle() {
	var (
		res []byte
	)
	for {
		bs = <-model.ByteStream
		res = append(res, bs...)
		if len(res) > 100000000 {
			list = bytes.Split(res, sep)
			filter(list)
			res = nil
		}
	}
}

var count int

func filter(list [][]byte) {
	st := time.Now()
	var res = false
	for _, v := range list {
		i := bytes.Index(v, sep2)
		if i == -1 {
			count++
			// if k > 0 {
			// 	fmt.Print(b2s(list[k-1]))
			// 	fmt.Println("   ",b2s(list[k]))
			// }
			// fmt.Println(list[k-1])
			continue
		}
		fspan.Tid = b2s(v[:i])
		fspan.Data = b2s(v)
		model.Stream <- fspan
		res = bytes.Contains(v, sep3)
		if res {
			ws.WriteTid(s2b(fspan.Tid))
			model.Mux.Lock()
			model.ErrTid[fspan.Tid] = ""
			model.Mux.Unlock()
			continue
		}
		res = bytes.Contains(v, sep4)
		if res {
			res = bytes.Contains(v, sep5)
			if !res {
				ws.WriteTid(s2b(fspan.Tid))
				model.Mux.Lock()
				model.ErrTid[fspan.Tid] = ""
				model.Mux.Unlock()
			}
		}
	}
	fmt.Println(time.Now().Sub(st))
	fmt.Println(count)
}
