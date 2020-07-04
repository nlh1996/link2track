package utils

import (
	"cloud/env"
	"cloud/model"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unsafe"
)

// Md5 .
func Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HTTPPost .
func HTTPPost() {
	DataURLVal := url.Values{}
	mjson, _ := json.Marshal(model.Result)
	mString := Bytes2str(mjson)
	DataURLVal.Add("result", mString)
	resp, err := http.Post(env.URL,
		"application/x-www-form-urlencoded",
		strings.NewReader(DataURLVal.Encode()))
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

// Str2bytes .
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2str .
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
