package main

import (
	"happy-dog/model"
	"happy-dog/routes"
)

func main() {
	model.InitDb()
	model.InitRedis()
	model.InitConn()
	routes.InitRouter()
}
