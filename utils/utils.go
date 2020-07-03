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
	mString := string(mjson)
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
