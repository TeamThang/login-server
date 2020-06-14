// 提醒模块
package notify

import (
	"github.com/gin-gonic/gin"
	"login-server/model/base"
	"login-server/model/db/postgre"
	"login-server/model/log"
)

// 添加提醒邮箱
func AddEmailHandler(c *gin.Context) {
	var email postgre.Email
	c.BindJSON(&email)
	err := AddEmail(email.UserID, email.Email, email.Subscribed)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 添加提醒邮箱
func DelEmailHandler(c *gin.Context) {
	emailID, err := base.GetUriID("EmailID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	err = DelEmail(emailID)
	if err == nil {
		base.RespondWithCode(204, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 查询用户的提醒邮箱
func QueryUserEmailHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	emails, err := QueryAllEmail(userID)
	data := map[string]interface{}{"data": emails}
	if err == nil {
		base.RespondWithCode(201, "", data, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), data, c)
	}
}

// 修改提醒邮箱为订阅
func SubEmailHandler(c *gin.Context) {
	emailID, err := base.GetUriID("EmailID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	err = SubEmail(emailID)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 修改提醒邮箱为非订阅
func UnSubEmailHandler(c *gin.Context) {
	emailID, err := base.GetUriID("EmailID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	err = UnSubEmail(emailID)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 邮件提醒发送请求
type ReqSendEmail struct {
	Subject     string
	ContentType string
	Content     string
}

// 发送邮件提醒
func SendEmailHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	var reqSendEmail ReqSendEmail
	c.BindJSON(&reqSendEmail)
	log.Debug("ReqSendEmail: %v\n", &reqSendEmail)
	err = SendNotifyEmail(userID, reqSendEmail.Subject, reqSendEmail.ContentType, reqSendEmail.Content)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}
