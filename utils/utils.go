package utils

import (
	"crypto/md5"
	"fmt"
)

// Md5 .
func Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}
