package utils

import (
	"cloud/env"
	"cloud/model"
	"cloud/socket"
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	buffer []byte
	temp   string
	list   []string
	length int
	start  time.Time
	end    time.Time
)

func init() {
	buffer = make([]byte, 30000000)
}

// Md5 .
func Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetRes .
func GetRes(url string) {
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
	start = time.Now()

	readData(resp)
}

func readData(resp *http.Response) {
	for {
		n, err := resp.Body.Read(buffer)
		if n == 0 || err != nil {
			model.EndSign = 1
			end = time.Now()
			fmt.Println("读取结束", end.Sub(start), n, err)
			// resp.Body.Close()
			break
		}
		list = strings.Split(string(buffer[:n]), "\n")
		length = len(list)
		temp = list[length-1]
		list[0] = temp + list[0]
		filter(list[:length-1])
	}
}

func filter(list []string) {
	var res = false
	span := model.Span{}
	for _, v := range list {
		arr := strings.Split(v, "|")
		span.Tid = arr[0]
		span.Data = v + "\n"
		model.Stream <- span
		if len(arr) < 9 {
			continue
		}
		res = strings.Contains(arr[8], "error=1")
		if res {
			socket.Write(arr[0])
			continue
		}
		res = strings.Contains(arr[8], "code=4")
		if res {
			socket.Write(arr[0])
			continue
		}
		res = strings.Contains(arr[8], "code=5")
		if res {
			socket.Write(arr[0])
		}
	}
}
