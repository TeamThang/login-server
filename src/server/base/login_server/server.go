package login_server

import (
	"io/ioutil"
	"encoding/json"
	golog "log"

	"server/util"
	"server/log"
	"server/conf"
	"server/db/postgre"
	"server/db/redis"
)

var HttpAdr = ""

func InitServer(confPath string, serverCfg string) {
	data, err := ioutil.ReadFile(util.PathJoin(confPath, serverCfg))
	if err != nil {
		log.Fatal("%v \n", err)
	}
	err = json.Unmarshal(data, &conf.Server)
	if err != nil {
		log.Fatal("%v \n", err)
	}
	HttpAdr = conf.Server.HttpAddr
	// 初始化日志 logger
	if conf.Server.LogLevel != "" { //日志级别不为空
		logger, err := log.New(conf.Server.LogLevel, conf.Server.LogPath, golog.LstdFlags) //创建一个logger
		if err != nil {
			panic(err)
		}
		log.Export(logger) //替换默认的gLogger
	}
	log.Release("log level: %v \n", conf.Server.LogLevel)
	conf.InitConfig(confPath)
	// 初始化postgre
	postgre.InitDB()
	// 初始化redis
	redis.InitPool()
}
