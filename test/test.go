package main

import (
	"cloud/env"
	"cloud/model"
	"cloud/server"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// var buffer []byte

func main() {
	// env.Client = &http.Client{
	// 	Timeout: time.Second * 60,
	// }
	// url := "http://www.yinghuo2018.com/download/trace1.data"
	// start := time.Now()
	// get(url, 0, 1000000)
	// end := time.Now()
	// fmt.Println(end.Sub(start))

	// start = time.Now()
	// getRes(url, 0, 1000000)
	// end = time.Now()
	// fmt.Println(end.Sub(start))
	model.Init()
	// 开启socket服务端
	server.Server()
}

// func getRes(url string, i int, size int) string {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Fatalf("Invalid url for downloading: %s, error: %v", url, err)
// 	}
// 	s := fmt.Sprintf("bytes=%d-%d", i*size, (i+1)*size)
// 	req.Header.Set("Range", s)
// 	req.Header.Set("Accept-Charset", "utf-8")
// 	resp, err := env.Client.Do(req)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
// 	var res string
// 	for {
// 		n, err := resp.Body.Read(buffer)
// 		if err != nil || n == 0 {
// 			fmt.Println("出现错误")
// 			break
// 		}
// 		res += string(buffer[:n])
// 	}
// 	if resp != nil {
// 		resp.Body.Close()
// 	}
// 	return res
// }

// func getRes2(url string, i int, size int) string {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Fatalf("Invalid url for downloading: %s, error: %v", url, err)
// 	}
// 	s := fmt.Sprintf("bytes=%d-%d", i*size, (i+1)*size)
// 	req.Header.Set("Range", s)
// 	req.Header.Set("Accept-Charset", "utf-8")
// 	resp, err := env.Client.Do(req)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	if resp != nil {
// 		resp.Body.Close()
// 	}
// 	return string(body)
// }

// func testStrings() {
// 	var str = "Hello world!!!"
// 	var res = false
// 	var msg string
// 	reg := regexp.MustCompile(`He|ld`)
// 	start := time.Now()
// 	for i := 0; i < 1000000; i++ {
// 		res = strings.Contains(str, "He")
// 		if !res {
// 			res = strings.Contains(str, "ld")
// 		}
// 	}
// 	end := time.Now()
// 	fmt.Println(end.Sub(start))

// 	start = time.Now()
// 	for i := 0; i < 1000000; i++ {
// 		msg = reg.FindString(str)
// 		if msg != "" {

// 		}
// 	}
// 	end = time.Now()
// 	fmt.Println(end.Sub(start))
// }

// // func testRouter() {
// // 	env.Port = "8000"
// // 	go router.Init()
// // 	<-env.ResStart
// // 	fmt.Println(111)
// // 	select {}
// // }

// func httpRead() {
// 	// buf := make([]byte, 1000000)
// 	// for {
// 	// 	n, err := resp.Body.Read(buf)
// 	// 	if err != nil || n == 0 {
// 	// 		fmt.Println("出现错误")
// 	// 		break
// 	// 	}
// 	// 	result += string(buf[:n])
// 	// }
// }

// func testsort() {
// 	list := model.Spans{}
// 	span := model.Span{}
// 	span.Time = "1587457762873000"
// 	list = append(list, span)
// 	span.Time = "1587457762872000"
// 	list = append(list, span)
// 	span.Time = "1587457762874000"
// 	list = append(list, span)
// 	fmt.Println(list)
// 	sort.Sort(list)
// 	fmt.Println(list)
// }

func get(url string, i int, size int) string {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	s := fmt.Sprintf("bytes=%d-%d", i*size, (i+1)*size)
	req.Header.Set("Range", s)
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Accept-Encoding", "gzip")
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	if err := fasthttp.DoTimeout(req, resp, time.Second*30); err != nil {
		log.Println(err)
	}
	return string(resp.Body())
}

func getRes(url string, i int, size int) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Invalid url for downloading: %s, error: %v", url, err)
	}
	s := fmt.Sprintf("bytes=%d-%d", i*size, (i+1)*size)
	req.Header.Set("Range", s)
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := env.Client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(body)
}
