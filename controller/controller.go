package controller

import (
	"cloud/env"
	"cloud/model"
	"cloud/socket"
	"cloud/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// Ready .
func Ready(c *gin.Context) {
	c.String(200, "ok")
}

var index int
// SetParameter .
func SetParameter(c *gin.Context) {
	index ++
	env.ResPort = c.Query("port")
	log.Println(env.ResPort)
	if env.Port != "8002" && index == 1 {
		go start()
	}
	c.String(200, "ok")
}

func start() {
	if env.Port == "8000" {
		env.URL = "http://localhost:" + env.ResPort + "/trace1.data"
	}
	if env.Port == "8001" {
		env.URL = "http://localhost:" + env.ResPort + "/trace2.data"
	}

	go streamHandle()

	utils.GetRes(env.URL)
}

func streamHandle() {
	for {
		if model.EndSign == 1 {
			for {
				span := <-model.Stream
				model.Mux.Lock()
				_, ok := model.ErrTid[span.Tid]
				model.Mux.Unlock()
				if ok {
					socket.Write(span.Data)
				}
				if len(model.Stream) == 0 {
					socket.Write("end")
					return
				}
			}
		}
		if len(model.Stream) > 30000 {
			span := <-model.Stream
			model.Mux.Lock()
			_, ok := model.ErrTid[span.Tid]
			model.Mux.Unlock()
			if ok {
				socket.Write(span.Data)
			}
		}
	}
}
