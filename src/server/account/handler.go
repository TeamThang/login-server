// 账户增删改查handler
// 解析上报参数
// 调用github.com/login_server/account中函数
package account

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"server/log"
	"server/base"
)

// 用户注册handler
// 解析json数据，并调用创建方法
func RegisterHandler(c *gin.Context) {
	var userCreate UserCreate
	c.BindJSON(&userCreate)
	userID, err := CreateUser(&userCreate)
	userIDMap := map[string]interface{}{"userID": userID}
	if err == nil {
		base.RespondWithCode(201, "", userIDMap, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), userIDMap, c)
	}
}

// 用户销户handler
// 解析json数据，并调用创建方法
func UnRegisterHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	err = DeleteUser(uint(userID))
	if err == nil {
		base.RespondWithCode(204, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 用户修改请求
type UserUpdate struct {
	BasicInfo json.RawMessage
	UserInfo  json.RawMessage
}

// 用户信息更新handler
func UpdateHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	var userUpdate UserUpdate
	c.BindJSON(&userUpdate)
	err = UpdateUser(uint(userID), &UserCreate{BasicInfo: userUpdate.BasicInfo, UserInfo: userUpdate.UserInfo})
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 修改密码请求
type ChanPasswd struct {
	OldPW string
	NewPW string
}

// 用户修改密码
// 需要验证原始密码
func ChangePasswordHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	var chanPasswd ChanPasswd
	c.BindJSON(&chanPasswd)
	log.Debug("ChanPasswd: %v\n", chanPasswd)
	err = ChangPassword(uint(userID), chanPasswd.OldPW, chanPasswd.NewPW)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 根据UserID查询用户信息
func QueryOneHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	var userQuery *UserQuery
	userQuery, err = QueryOneUser(uint(userID))
	data := map[string]interface{}{"data": userQuery}
	if err == nil {
		base.RespondWithCode(200, "", data, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), data, c)
	}
}

// 邮箱找回密码
func EmailRestPwdSendHandler(c *gin.Context) {
	var restPd RestPassword
	c.BindJSON(&restPd)
	userID, err := EmailFindPassword(&restPd)
	data := map[string]interface{}{"userID": userID}
	if err == nil {
		base.RespondWithCode(200, "", data, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), data, c)
	}
}

// 邮箱重置密码
func EmailRestPwdCheckHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	verifyCode, err := base.GetUriStr("Code", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	var chanPasswd ChanPasswd
	c.BindJSON(&chanPasswd)
	log.Debug("UserID: %d, Code: %s, ChanPasswd: %v\n", userID, verifyCode, chanPasswd)
	err = EmailResetPassword(userID, verifyCode, chanPasswd.NewPW)
	if err == nil {
		base.RespondWithCode(200, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}
