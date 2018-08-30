// 发送验证码并验证账户
package verify

import (
	"github.com/gin-gonic/gin"
	"server/base"
)

// 发送邮件验证码
func SendEmailHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	err = SendVerifyEmail(userID)
	if err == nil {
		base.RespondWithCode(200, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 验证邮件验证码
func CheckEmailHandler(c *gin.Context) {
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	verifyCode, err := base.GetUriStr("Code", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	res, err := CheckVerifyCode(userID, verifyCode)
	data := map[string]interface{}{"result": res}
	if err == nil {
		base.RespondWithCode(201, "", data, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), data, c)
	}
}
