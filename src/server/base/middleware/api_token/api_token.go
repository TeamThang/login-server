package api_token

import (
	"github.com/gin-gonic/gin"

	"server/login/token"
	"server/base"
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

// token鉴权中间件
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenID string
		sessionID, err := c.Cookie("session_id") // 优先取session_id
		if err != nil {
			tokenID = c.GetHeader("X-Request-Id") // 其次取请求头的X-Request-Id
		} else {
			tokenID = sessionID
		}
		if tokenID == "" {
			base.RespondWithCode(401, "API token required", nil, c)
			return
		}
		authFlag, err := token.TokenFlagAuth(tokenID)
		if err != nil {
			base.RespondWithCode(401, fmt.Sprintf("auth check failed err: %v", err), nil, c)
			return
		}
		if !authFlag {
			base.RespondWithCode(401, "insufficient privileges", nil, c)
			return
		}
		c.Next()
	}
}

type ReqID struct {
	UserID uint
}

// user鉴权中间件
func UserAuthUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenID string
		sessionID, err := c.Cookie("session_id") // 优先取session_id
		if err != nil {
			tokenID = c.GetHeader("X-Request-Id") // 其次取请求头的X-Request-Id
		} else {
			tokenID = sessionID
		}
		if tokenID == "" {
			base.RespondWithCode(401, "API token required", nil, c)
			return
		}
		authAdminFlag, err := token.TokenFlagAuth(tokenID) // 具有登陆服务admin权限，则不需要UserID和token匹配鉴权
		if err == nil {
			if authAdminFlag {
				c.Next()
				return
			}
		}
		// 从uri或请求中获取UserID
		userID, err := base.GetUriID("UserID", c)
		if err != nil {
			var reqID ReqID
			contents, _ := ioutil.ReadAll(c.Request.Body)                // 读取body中数据流
			json.Unmarshal(contents, &reqID)                             // 检测UserID
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(contents)) // 重新写入Body

			userID = reqID.UserID
			if userID == 0 {
				base.RespondWithCode(401, fmt.Sprintf("user id is needed"), nil, c)
				return
			}
		}

		authFlag, err := token.TokenUserAuth(tokenID, userID)
		if err != nil {
			base.RespondWithCode(401, fmt.Sprintf("auth check failed err: %v", err), nil, c)
			return
		}
		if !authFlag {
			base.RespondWithCode(401, "user id is not matched", nil, c)
			return
		}
		c.Next()
	}
}
