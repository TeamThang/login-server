package redis

import (
	"fmt"
	"time"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"server/conf"
	"server/log"
)

var RedisClient *redis.Pool

// 初始化redis连接池
func InitPool() {
	config := conf.Config.Redis
	// 建立连接池
	RedisClient = &redis.Pool{
		MaxIdle:     1,                 // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		MaxActive:   10,                // 最大的激活连接数，表示同时最多有N个连接
		IdleTimeout: 10 * time.Second, // 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) {
			redis_tcp := fmt.Sprintf("%s:%s", config.Host, strconv.Itoa(config.Port))
			log.Debug("redis tcp: %v \n", redis_tcp)
			c, err := redis.Dial("tcp", redis_tcp)
			if err != nil {
				return nil, err
			}
			// 验证密码，如果有密码
			if config.PassWord != "" {
				if _, err := c.Do("AUTH", config.PassWord); err != nil {
					c.Close()
					return nil, err
				}
			}
			// 选择db
			c.Do("SELECT", config.DB)
			return c, nil
		},
	}
}

// 使用连接池的Do
func Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	rc := RedisClient.Get()
	defer rc.Close()
	res, err := rc.Do(commandName, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
