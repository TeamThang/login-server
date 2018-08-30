package main

import (
	"github.com/gin-gonic/gin"
	"server/log"
	"server/base/router"
	"server/base/login_server"
)

var confPath = "/Users/yzy/Develop/job_bitmain/login_server/conf/"

func init() {
	login_server.InitServer(confPath, "server.json")
}

func main() {
	server := gin.New()
	server.Use(gin.Logger())
	router.SetRouter(server)
	server.Run(login_server.HttpAdr)
	log.Close() // 关闭日志
}
