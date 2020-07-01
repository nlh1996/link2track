package env

import (
	"net/http"
	"time"
)

var (
	Port    string
	ResPort string
	Client  *http.Client
	URL     string
)

func init() {
	Client = &http.Client{
		Timeout: time.Second * 30,
	}
}
