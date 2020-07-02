package env

import (
	"net/http"
	"time"
)

var (
	Port       string
	ResPort    string
	Client     *http.Client
	URL        string
	BufferSize int
	StreamSize int
)

func init() {
	Client = &http.Client{
		Timeout: time.Second * 30,
	}
	BufferSize = 1000
	StreamSize = 500100
}
