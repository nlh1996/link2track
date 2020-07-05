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
		Timeout: time.Second * 400,
	}

	if os.Getenv("SERVER_PORT") == "" {
		BufferSize = 100000
		StreamSize = 1001000
	} else {
		BufferSize = 10000000
		StreamSize = 2001000
	}
}
