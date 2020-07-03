package main

import (
	"cloud/env"
	"cloud/model"
	"cloud/router"
	"cloud/server"
	"cloud/socket"
	"os"
	"runtime"
)

func initData() {
	env.Port = os.Getenv("SERVER_PORT")
	if env.Port == "8000" {
		runtime.GOMAXPROCS(4)
		env.URL = "http://www.yinghuo2018.com/download/trace1.data"
	}
	if env.Port == "8001" {
		runtime.GOMAXPROCS(4)
		env.URL = "http://www.yinghuo2018.com/download/trace2.data"
	}
	if env.Port == "8002" {
		runtime.GOMAXPROCS(2)
	}
	if env.Port == "" {
		// runtime.GOMAXPROCS(4)
		// env.URL = "http://192.168.0.4/trace1.data"
		env.URL = "http://www.yinghuo2018.com/download/trace1.data"
		env.Port = "8000"
	}
}

func main() {
	// 数据初始化
	initData()
	model.Init()

	// backend
	if env.Port == "8002" {
		// 开启http服务
		go router.Init()
		// 开启socket服务端
		go server.Server(":8003")
		server.Server(":8004")
	}

	// 开启socket客户端
	socket.Init()
	// 开启http服务
	go router.Init()

	select {}
}
