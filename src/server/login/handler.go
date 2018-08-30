package login

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"server/log"
	"server/login/login"
	"server/account"
	"server/base"
)

// 登录
func LoginHandler(c *gin.Context) {
	loginInfo := &login.LoginInfo{}
	c.BindJSON(loginInfo)
	log.Debug("loginInfo: %v\n", *loginInfo)
	sessionID, userID, loginName, err := login.Login(loginInfo)
	data := map[string]interface{}{"userID": userID, "loginName": loginName}
	if err == nil {
		c.SetCookie("session_id", sessionID, int(loginInfo.Duration), "", "", true, true)
		c.Writer.Header().Set("X-Request-Id", sessionID)
		c.SetCookie("server", loginInfo.Right.Server, int(loginInfo.Duration), "", "", true, true)
		base.RespondWithCode(201, "login success", data, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), data, c)
	}
}

// 注销
func LogoutHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		base.RespondWithCode(400, "unlogin", nil, c)
	}
	err = login.Logout(sessionID)
	if err == nil {
		base.RespondWithCode(200, "logout success", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 获取登录用户信息
func GetUserInfoHandler(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		base.RespondWithCode(400, fmt.Sprintf("Get session_id from cookie err: %s", err.Error()), nil, c)
	}
	var userQuery *account.UserQuery
	userQuery, err = login.GetUserInfo(sessionID)
	data := map[string]interface{}{"data": userQuery}
	if err == nil {
		base.RespondWithCode(200, "", data, c)
	} else {
		base.RespondWithCode(400, err.Error(), data, c)
	}
}
