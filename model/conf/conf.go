package conf

import (
	"encoding/json"
	"io/ioutil"

	"login-server/model/log"
	"login-server/util"
)

type Postgre struct {
	Host     string
	Port     int
	User     string
	DbName   string
	PassWord string
}

type Redis struct {
	Host     string
	Port     int
	PassWord string
	DB       int
}

type Email struct {
	Host     string
	Port     int
	User     string
	PassWord string
	Sender   string
}

var Config struct {
	Postgre Postgre
	Redis   Redis
	Email   Email
}

func InitConfig(confPath string) {
	data, err := ioutil.ReadFile(util.PathJoin(confPath, "config.json"))
	if err != nil {
		log.Fatal("%v \n", err)
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatal("%v \n", err)
	}
	log.Release("Config: %v \n", Config)
}
