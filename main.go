package main

import (
	"cloud/env"
	"cloud/model"
	"cloud/router"
	"cloud/ws"
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
		env.URL = "http://192.168.0.4/trace2.data"
		//env.URL = "http://www.yinghuo2018.com/download/trace1.data"
		env.Port = "8000"
	}
}

func main() {
	// 数据初始化
	initData()
	model.Init()

	if env.Port == "8002" {
		router.Init()
	}

	//连接websocket服务
	ws.Dial()
	// 开启http服务
	go router.Init()

	select {}
}
