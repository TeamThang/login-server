package token

import (
	"fmt"
	rredis "github.com/gomodule/redigo/redis"
	"strconv"
	"strings"

	"login-server/model/db/redis"
	"login-server/model/log"
	"login-server/util"
)

// token采用k-v存储redis中，设置过期时间
// token为uuid, 值为User.ID
const TokenKeyFmt string = "LOGIN_TOKEN:%s" // 存入redis的key格式
const TokenValFmt string = "%s_%s"          // 存入redis的value格式，， {userID}_{admin flag}
const TokenDuration uint = 36000            // 默认token持续时间，单位:s

// 设置session到redis
// loginName: 登录名; reqRight: 登录鉴权信息; duration: session超时时间
func SetSessionID(userID uint, duration uint, login_admin uint) (string, error) {
	token, err := genToken()
	if err != nil {
		return "", fmt.Errorf("token generate err: %s", err)
	}
	tokenVal := fmt.Sprintf(TokenValFmt, strconv.Itoa(int(userID)), strconv.Itoa(int(login_admin)))
	log.Debug("token: %s, val: %s", token, tokenVal)
	tokenKey := fmt.Sprintf(TokenKeyFmt, token) // 存入redis时，格式化格式为 TOKEN:uuid
	_, err = redis.Do("set", tokenKey, tokenVal)
	if err != nil {
		return "", fmt.Errorf("token set 1 err: %s", err)
	}
	_, err = redis.Do("expire", tokenKey, duration)
	if err != nil {
		return "", fmt.Errorf("token set 2 err: %s", err)
	}
	return token, nil
}

// 从redis中读取token的值
// 获取对应的userID
func GetTokenValByID(token string) (uint, uint, error) {
	tokenKey := fmt.Sprintf(TokenKeyFmt, token) // 存入redis时，格式化格式为 TOKEN:uuid
	res, err := rredis.String(redis.Do("get", tokenKey))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			return 0, 0, fmt.Errorf("unlogin")
		}
		return 0, 0, fmt.Errorf("token val get by id err: %s", err)
	}
	value := strings.Split(res, "_")
	userID, err := strconv.Atoi(value[0])
	if err != nil {
		return 0, 0, fmt.Errorf("token val get from id err: %s", err)
	}
	authFlag, err := strconv.Atoi(value[1])
	if err != nil {
		return 0, 0, fmt.Errorf("token val get from id err: %s", err)
	}
	return uint(userID), uint(authFlag), nil
}

// 从redis删除对应的session
func DelSessionID(token string) error {
	tokenKey := fmt.Sprintf(TokenKeyFmt, token) // 存入redis时，格式化格式为 TOKEN:uuid
	_, err := rredis.String(redis.Do("get", tokenKey))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			return fmt.Errorf("unlogin")
		}
		return err
	}
	_, err = redis.Do("del", tokenKey)
	if err != nil {
		return fmt.Errorf("token delete err: %s", err)
	}
	return nil
}

// 从redis清除对应userID的session
func CleanSessionID(userID uint) error {
	sessionQueryKey := fmt.Sprintf(TokenKeyFmt, "*")
	tokens, err := rredis.Strings(redis.Do("keys", sessionQueryKey))
	if err != nil {
		return fmt.Errorf("token clean err: %s", err)
	}
	for _, tokenKey := range tokens {
		res, err := rredis.String(redis.Do("get", tokenKey))
		if err != nil {
			return fmt.Errorf("token clean err:  %s", err)
		}
		value := strings.Split(res, "_")
		uID, err := strconv.Atoi(value[0])
		if err != nil {
			return fmt.Errorf("token val get from id err: %s", err)
		}
		if uint(uID) == userID {
			_, err := redis.Do("del", tokenKey)
			if err != nil {
				return fmt.Errorf("token clean err:  %s", err)
			}
		}
	}
	return nil
}

// 生成token
func genToken() (string, error) {
	uuid := util.GetUUID()
	return uuid, nil
}

// 根据token判断是否具有操作login_server权限
func TokenFlagAuth(token string) (bool, error) {
	_, AdminFlag, err := GetTokenValByID(token)
	if err != nil {
		return false, err
	}
	if AdminFlag == 1 {
		return true, err
	}
	return false, err
}

// 根据token判断是否具有修改查看权限
func TokenUserAuth(token string, userID uint) (bool, error) {
	uID, _, err := GetTokenValByID(token)
	if err != nil {
		return false, err
	}
	if uID == userID {
		return true, err
	}
	return false, err
}
