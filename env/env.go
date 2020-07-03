package env

import (
	"net/http"
	"os"
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
		Timeout: time.Second * 120,
	}
	if os.Getenv("SERVER_PORT") == "" {
		BufferSize = 1000
	} else {
		BufferSize = 1000000
	}
	StreamSize = 500100
}
