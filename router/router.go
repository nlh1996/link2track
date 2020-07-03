package router

import (
	"cloud/controller"
	"cloud/env"
	"cloud/middleware"
	"cloud/ws"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Init .
func Init() {
	router := gin.Default()
	router.Use(middleware.CrossDomain())
	router.GET("/", Handler)
	router.GET("/ready", controller.Ready)
	router.GET("/setParameter", controller.SetParameter)
	router.Run(":" + env.Port)
}

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	router *ws.Router
)

func init() {
	router = ws.NewRouter()
	// router.AddRouter(100, scene.SayHello)
}

// Handler .
func Handler(c *gin.Context) {
	// 升级get请求为webSocket长连接
	connection, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := ws.NewConnection(connection, router)
	if err != nil {
		connection.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	conn.Start()
}
