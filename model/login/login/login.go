// 登录登出
package login

import (
	"fmt"
	"time"

	"login-server/model/account"
	"login-server/model/db/postgre"
	"login-server/model/log"
	"login-server/model/login/token"
	"login-server/util"
)

type Right struct {
	Server string
	Name   string
}

type LoginInfo struct {
	LoginName string
	Password  string
	Mobile    string
	Email     string
	Right     Right
	Duration  uint
}

// 登录
func Login(loginInfo *LoginInfo) (sessionID string, userID uint, loginName string, err error) {
	if loginInfo.Password == "" {
		err = fmt.Errorf("Password is needed")
		return
	}
	if loginInfo.LoginName == "" && loginInfo.Mobile == "" && loginInfo.Email == "" {
		err = fmt.Errorf(("Please use LoginName or Mobile or Email to login"))
		return
	}
	var user *postgre.User
	if loginInfo.LoginName != "" { // 优先登录名登录
		sessionID, user, err = loginByLoginName(loginInfo)
		if err != nil {
			return
		}
	}
	if loginInfo.Mobile != "" { // 其次手机登录
		sessionID, user, err = loginByMobile(loginInfo)
		if err != nil {
			return
		}
	}
	if loginInfo.Email != "" { // 最后邮箱登录
		sessionID, user, err = loginByEmail(loginInfo)
		if err != nil {
			return
		}
	}
	if user != nil {
		userID = user.ID
		loginName = user.LoginName
	}
	return
}

// 通过登录名登录
func loginByLoginName(loginInfo *LoginInfo) (string, *postgre.User, error) {
	user := &postgre.User{}
	db := postgre.DB.Where(&postgre.User{LoginName: loginInfo.LoginName}).First(&user)
	if db.Error != nil {
		if db.RecordNotFound() {
			return "", nil, fmt.Errorf("LoginName is not existed")
		}
		return "", nil, db.Error
	}
	sessionID, err := checkAndLogin(user, loginInfo)
	if err != nil {
		return "", nil, err
	}
	return sessionID, user, nil
}

// 通过手机登录
func loginByMobile(loginInfo *LoginInfo) (string, *postgre.User, error) {
	user := &postgre.User{}
	db := postgre.DB.Where(&postgre.User{Mobile: loginInfo.Mobile}).First(&user)
	if db.Error != nil {
		if string(db.Error.Error()) == "record not found" {
			return "", nil, fmt.Errorf("Mobile is not existed")
		}
		return "", nil, db.Error
	}
	sessionID, err := checkAndLogin(user, loginInfo)
	if err != nil {
		return "", nil, err
	}
	return sessionID, user, nil
}

// 通过邮箱登录
func loginByEmail(loginInfo *LoginInfo) (string, *postgre.User, error) {
	user := &postgre.User{}
	db := postgre.DB.Where(&postgre.User{Email: loginInfo.Email}).First(&user)
	if db.Error != nil {
		if db.RecordNotFound() {
			return "", nil, fmt.Errorf("Email is not existed")
		}
		return "", nil, db.Error
	}
	sessionID, err := checkAndLogin(user, loginInfo)
	if err != nil {
		return "", nil, err
	}
	return sessionID, user, nil
}

// 鉴权并登录
// login_server权限管理: 1: admin(all,login)
func checkAndLogin(user *postgre.User, loginInfo *LoginInfo) (string, error) {
	if loginInfo.Duration == 0 {
		loginInfo.Duration = token.TokenDuration // 如果未传入持续时间，取默认时间
	}
	pd_md5 := util.GetMD5(loginInfo.Password)
	if pd_md5 != user.Password {
		return "", fmt.Errorf("password is not right")
	}
	err := checkRight(user.ID, &loginInfo.Right)
	if err != nil {
		return "", err
	}
	adminFlag, err := checkLoginRight(user.ID) // 检查是否具有login_server的admin权限
	if err != nil {
		return "", err
	}

	sessionID, err := token.SetSessionID(user.ID, loginInfo.Duration, adminFlag)
	if err != nil {
		return "", err
	}
	go updateUser(user)
	go writeLoginLog(user.ID, "login", loginInfo.Duration, sessionID)
	return sessionID, nil
}

// 注销
func Logout(sessionID string) error {
	err := token.DelSessionID(sessionID)
	if err != nil {
		return err
	}
	go updateLoginLog(sessionID, "login")
	return nil
}

// 获得用户信息
func GetUserInfo(sessionID string) (*account.UserQuery, error) {
	userID, _, err := token.GetTokenValByID(sessionID)
	if err != nil {
		return nil, err
	}
	userQuery := &account.UserQuery{}
	userQuery, err = account.QueryOneUser(userID)
	if err != nil {
		return nil, err
	}
	return userQuery, nil
}

// 登录后更新User表
func updateUser(user *postgre.User) {
	db := postgre.DB.Model(&user).Updates(&postgre.User{LoginTime: time.Now(), LoginCount: user.LoginCount + 1})
	if db.Error != nil {
		log.Error("update user failed, userID: %v", user.ID)
	}
}

// 登录后记录LoginLog表
func writeLoginLog(userID uint, opType string, duration uint, token string) {
	db := postgre.DB.Create(&postgre.LoginLog{UserID: userID, OpType: opType, OpTime: time.Now(), Duration: duration, Token: token})
	if db.Error != nil {
		log.Error("write login_log failed, userID: %v", userID)
	}
}

// 注销事更新LoginLog表
// 记录登录到注销的时间
func updateLoginLog(token string, opType string) {
	res := &postgre.LoginLog{}
	db := postgre.DB.Where(&postgre.LoginLog{Token: token, OpType: opType}).First(&res)
	if db.Error != nil {
		log.Error("update login_log failed, token: %s", token)
	}
	time_delta := uint(time.Now().Sub(res.OpTime).Seconds())
	db = postgre.DB.Model(&res).Updates(&postgre.LoginLog{Duration: time_delta})
	if db.Error != nil {
		log.Error("update login_log failed, token: %s", token)
	}
}
