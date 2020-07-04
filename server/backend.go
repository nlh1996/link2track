package main

import (
	"cloud/env"
	"cloud/model"
	"cloud/router"
)


func main() {
	model.SInit()
	env.Port = "8002"
	router.Init()
}
