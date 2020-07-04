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
	// buffer = make([]byte, env.BufferSize)
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

	//go streamHandle()
	start = time.Now()
	getTid()
	fmt.Println("第一次请求时间", time.Now().Sub(start))

	start = time.Now()
	getRes()
	fmt.Println("第二次请求时间", time.Now().Sub(start))
}

func getTid() {
	req, err := http.NewRequest("GET", env.URL, nil)
	if err != nil {
		log.Fatalf("Invalid url for downloading")
	}
	req.Header.Set("Accept-Charset", "utf-8")
	// var r string
	// if i == 9 {
	// 	r = fmt.Sprintf("%v-", 40000000*i)
	// } else {
	// 	r = fmt.Sprintf("%v-%v", 40000000*i, 40000000*(i+1))
	// }
	// req.Header.Set("Range", "bytes="+r)
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
			go filter(res)
			fmt.Println("读取结束", n, err)
			resp.Body.Close()
			return
		}
		if len(res) > 50000000 {
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

func readData2(resp *http.Response) {
	buffer := make([]byte, env.BufferSize)
	var res []byte
	for {
		n, err := resp.Body.Read(buffer)
		res = append(res, buffer[:n]...)
		if n == 0 || err != nil {
			go filter2(res, 1)
			fmt.Println("读取结束", n, err)
			//resp.Body.Close()
			return
		}
		if len(res) > 50000000 {
			go filter2(res, 0)
			res = nil
		}
	}
	// if body, err := ioutil.ReadAll(resp.Body); err != nil {
	// 	log.Println(err)
	// } else {
	// 	go filter(body)
	// }
}

func filter(bs []byte) {
	st := time.Now()
	var res = false
	list := strings.Split(b2s(bs), sep)
	for _, v := range list {
		arr := strings.Split(v, sep2)
		if len(arr) < 9 || len(arr[0]) < 14 {
			continue
		}
		res = strings.Contains(arr[8], sep3)
		if res {
			ws.WriteTid(s2b(arr[0]))
			continue
		}
		res = strings.Contains(arr[8], sep4)
		if res {
			res = strings.Contains(arr[8], sep5)
			if !res {
				ws.WriteTid(s2b(arr[0]))
			}
		}
	}
	fmt.Println("计算用时", time.Now().Sub(st))
}

func filter2(bs []byte, i int) {
	list := strings.Split(b2s(bs), sep)
	for _, v := range list {
		arr := strings.Split(v, sep2)
		if len(arr) < 9 || len(arr[0]) < 14 {
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
}