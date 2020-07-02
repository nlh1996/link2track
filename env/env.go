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
)

func init() {
	Client = &http.Client{
		Timeout: time.Second * 10,
	}
	BufferSize = 30000000
}
