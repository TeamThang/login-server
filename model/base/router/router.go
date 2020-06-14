package router

import (
	"github.com/gin-gonic/gin"
	"login-server/model/account"
	"login-server/model/authen"
	"login-server/model/base/middleware/api_token"
	"login-server/model/login"
	"login-server/model/notify"
	"login-server/model/verify"
)

func SetRouter(router *gin.Engine) {
	accountUrl := router.Group("/v1/account")
	{
		accountUrl.POST("", account.RegisterHandler)
		accountUrl.DELETE(":UserID", api_token.UserAuthUserMiddleware(), account.UnRegisterHandler)
		accountUrl.PUT(":UserID", api_token.UserAuthUserMiddleware(), account.UpdateHandler)
		accountUrl.GET(":UserID", api_token.UserAuthUserMiddleware(), account.QueryOneHandler)
	}
	emailRestPdUrl := router.Group("/v1/pwd")
	{
		emailRestPdUrl.PUT("user/:UserID", api_token.UserAuthUserMiddleware(), account.ChangePasswordHandler)
		emailRestPdUrl.POST("", account.EmailRestPwdSendHandler)
		emailRestPdUrl.GET("user/:UserID/code/:Code", account.EmailRestPwdCheckHandler)
	}

	authenUrl := router.Group("/v1/right")
	authenUrl.Use(api_token.TokenAuthMiddleware())
	{
		authenUrl.POST("", authen.CreateHandler)
		authenUrl.DELETE(":RightID", authen.DeleteHandler)
		authenUrl.PUT(":RightID", authen.UpdateHandler)
		authenUrl.GET("", authen.QueryAllHandler)
	}
	bindAuthUrl := router.Group("/v1/bind_auth")
	bindAuthUrl.Use(api_token.UserAuthUserMiddleware())
	{
		bindAuthUrl.POST("", authen.BindRightHandler)
		bindAuthUrl.DELETE("user/:UserID/right/:RightID", authen.UnBindRightHandler)
		bindAuthUrl.GET("user/:UserID", authen.QueryBindRightHandler)
	}
	verifyEmailUrl := router.Group("/v1/verify/email")
	{
		verifyEmailUrl.GET("user/:UserID", verify.SendEmailHandler)
		verifyEmailUrl.PUT("user/:UserID/code/:Code", verify.CheckEmailHandler)
	}
	notifyUrl := router.Group("/v1/notify/email")
	notifyUrl.Use(api_token.UserAuthUserMiddleware())
	{
		notifyUrl.POST("", notify.AddEmailHandler)
		notifyUrl.DELETE(":EmailID", notify.DelEmailHandler)
		notifyUrl.PUT(":EmailID/user/:UserID/sub", notify.SubEmailHandler)
		notifyUrl.PUT(":EmailID/user/:UserID/unsub", notify.UnSubEmailHandler)
		notifyUrl.GET("user/:UserID", notify.QueryUserEmailHandler)
		notifyUrl.POST("user/:UserID", notify.SendEmailHandler)
	}
	v1 := router.Group("/v1")
	{
		v1.POST("/login", login.LoginHandler)
		v1.GET("/logout", login.LogoutHandler)
		v1.GET("/get_user_info", login.GetUserInfoHandler)
	}
}
