package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"login-server/model/base/login_server"
	"login-server/model/base/router"
	"login-server/model/log"
)

var confPath = "/Users/yzy/Develop/job_bitmain/login_server/conf/"

func init() {
	flag.Parse()
	login_server.InitServer(confPath, "server.json")
}

func main() {
	server := gin.New()
	server.Use(gin.Logger())
	router.SetRouter(server)
	server.Run(login_server.HttpAdr)
	log.Close() // 关闭日志
}
