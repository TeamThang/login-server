// 账户数据格式化和增删改查
// 接受server/account/handler传入参数
// 调用github.com/login_server/db/postgre中方法操作数据库
package account

import (
	"encoding/json"
	"fmt"
	rredis "github.com/gomodule/redigo/redis"

	"login-server/model/authen"
	"login-server/model/db/postgre"
	"login-server/model/db/redis"
	"login-server/model/log"
	"login-server/model/login/token"
	"login-server/model/notify"
	"login-server/util"
	"strconv"
)

type UserCreate struct {
	BasicInfo json.RawMessage
	UserInfo  json.RawMessage
}

// 创建用户
func CreateUser(userCreate *UserCreate) (uint, error) {
	user, userInfo, err := makeUserInfo(userCreate)
	if err != nil {
		return 0, err
	}
	needParam := []string{"LoginName", "Password"}
	err = util.CheckAttrExist(user, needParam)
	if err != nil {
		return 0, err
	}
	if user.Email != "" { // 如果上报邮箱，检查邮箱格式
		err := util.CheckEmailFormat(user.Email)
		if err != nil {
			return 0, fmt.Errorf("please use right email address")
		}
	}

	user.Password = util.GetMD5(user.Password) // 密码用md5保存
	newUserID, err := postgre.UserCreate(user, userInfo)
	if err != nil {
		return 0, err
	}
	go func() { // 如果开户上报的source可以匹配到权限表中的server，就增加对应权限
		rightID, err := authen.QueryRightByServer(user.Source)
		if err == nil {
			if err := authen.BindRight(user.ID, rightID); err != nil {
				log.Error("new user[%d] bind right[%d] failed", user.ID, rightID)
			}
		}
	}()
	return newUserID, nil
}

// 解析用户信息参数到结构体
func makeUserInfo(userCreate *UserCreate) (*postgre.User, *postgre.UserInfo, error) {
	var user postgre.User
	var userInfo postgre.UserInfo
	var userInfoStr UserInfoStr
	log.Debug("basicInfo: %s\n", userCreate.BasicInfo)
	json.Unmarshal(userCreate.BasicInfo, &user) // 绑定数据到postgre.User
	log.Debug("userInfo: %s\n", userCreate.UserInfo)
	if userCreate.UserInfo != nil && string(userCreate.UserInfo) != "" {
		json.Unmarshal(userCreate.UserInfo, &userInfoStr) // 绑定数据到postgre.UserInfoStr
	}
	err := ToUserInfo(&userInfo, &userInfoStr)
	if err != nil {
		return nil, nil, err
	}
	return &user, &userInfo, nil
}

// 删除用户
// 先删除外键表，最后删除User
func DeleteUser(userID uint) error {
	_, err := postgre.CheckUserExist(userID) // 验证是否存在
	if err != nil {
		return err
	}
	err = postgre.UserRightRelationsDelete(userID) // 删除权限关联表
	if err != nil {
		return err
	}
	err = postgre.VerifyDelete(userID) // 删除验证表
	if err != nil {
		return err
	}
	err = postgre.UserInfoDelete(userID) // 删除用户信息表
	if err != nil {
		return err
	}
	err = postgre.EmailDelete(userID) // 删除提醒邮箱表
	if err != nil {
		return err
	}
	err = postgre.UserDelete(userID) // 删除用户账户表
	if err != nil {
		return err
	}
	go token.CleanSessionID(userID) // 清理redis中sessionid
	return nil
}

// 修改用户信息
func UpdateUser(userID uint, userUpdate *UserCreate) error {
	user, userInfo, err := makeUserInfo(userUpdate)
	if err != nil {
		return err
	}
	user.Password = "" // 这里不支持密码修改
	err = user.Update(userID)
	if err != nil {
		return err
	}
	if userInfo != nil {
		err = userInfo.Update(userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// 修改密码
// userID: User表主键id; opd: 老密码; npd: 新密码
func ChangPassword(userID uint, opd string, npd string) error {
	if userID == 0 {
		return fmt.Errorf("UserID is needed")
	}
	if opd == "" {
		return fmt.Errorf("OldPW is needed")
	}
	if npd == "" {
		return fmt.Errorf("NewPW is needed")
	}
	user := &postgre.User{}
	user, err := user.QueryUserById(userID)
	if err != nil {
		return err
	}
	opd_md5 := util.GetMD5(opd)
	if opd_md5 != user.Password {
		return fmt.Errorf("old password is not right")
	}
	npd_md5 := util.GetMD5(npd)
	tUser := &postgre.User{Password: npd_md5}
	err = tUser.Update(userID)
	if err != nil {
		return err
	}
	return nil
}

type UserQuery struct {
	BasicInfo UserStr
	UserInfo  UserInfoStr
}

// 查询单个用户
func QueryOneUser(userID uint) (*UserQuery, error) {
	var userQuery UserQuery
	user := &postgre.User{}
	user, err := user.QueryUserById(userID)
	if err != nil {
		return nil, err
	}
	userStr := &UserStr{}
	ToUserStr(user, userStr)
	userQuery.BasicInfo = *userStr
	userInfo := &postgre.UserInfo{}
	userInfo, err = userInfo.QueryUserInfoByUserID(userID)
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return &userQuery, nil // 没有userinfo就直接返回
	}
	userInfoStr := &UserInfoStr{}
	ToUserInfoStr(userInfo, userInfoStr)
	userQuery.UserInfo = *userInfoStr
	return &userQuery, nil
}

const (
	redisKeyFmt string = "RESET_PASSWORD:%s"
	duration           = 300
)

// 找回密码结构体
type RestPassword struct {
	LoginName string
	Mobile    string
	Email     string
	ResetUrl  string
}

// 邮箱找回密码
func EmailFindPassword(restPd *RestPassword) (uint, error) {
	if restPd.ResetUrl == "" {
		return 0, fmt.Errorf("passwd reset url is needed")
	}
	if restPd.LoginName == "" && restPd.Mobile == "" && restPd.Email == "" {
		return 0, fmt.Errorf("Please use LoginName or Mobile or Email to find passwd")
	}
	user := &postgre.User{}
	if restPd.LoginName != "" { // 优先登录名登录
		db := postgre.DB.Where(&postgre.User{LoginName: restPd.LoginName}).First(&user)
		if db.Error != nil {
			if string(db.Error.Error()) == "record not found" {
				return 0, fmt.Errorf("LoginName is not existed")
			}
			return 0, db.Error
		}
	}
	if restPd.Mobile != "" { // 其次手机登录
		db := postgre.DB.Where(&postgre.User{Mobile: restPd.Mobile}).First(&user)
		if db.Error != nil {
			if db.RecordNotFound() {
				return 0, fmt.Errorf("Mobile is not existed")
			}
			return 0, db.Error
		}
	}
	if restPd.Email != "" { // 最后邮箱登录
		db := postgre.DB.Where(&postgre.User{Email: restPd.Email}).First(&user)
		if db.Error != nil {
			if string(db.Error.Error()) == "record not found" {
				return 0, fmt.Errorf("Email is not existed")
			}
			return 0, db.Error
		}
	}
	verifyCode := util.RandStringRunes(5)
	err := resetVerifyKey(user.ID, verifyCode)
	if err != nil {
		return 0, err
	}
	emails := []string{user.Email}
	restUrl := fmt.Sprintf("%s/user/%d/code/%s", restPd.ResetUrl, user.ID, verifyCode)
	subject := "Reset Password"
	contentType := "text/html"
	content := fmt.Sprintf("Please click url to reset your password:	%s\n in 5 min", restUrl)
	go notify.SendEmail(user.ID, subject, contentType, content, &emails) // 异步发送邮寄
	return user.ID, nil
}

// 邮箱重置密码
func EmailResetPassword(userID uint, verifyCode string, newPassword string) error {
	if userID == 0 {
		return fmt.Errorf("UserID is needed")
	}
	if verifyCode == "" {
		return fmt.Errorf("VerifyCode is needed")
	}
	vc, err := getVerifyKey(userID)
	if err != nil {
		return err
	}
	if vc != verifyCode {
		return fmt.Errorf("url is not valid")
	}
	npd_md5 := util.GetMD5(newPassword)
	tUser := &postgre.User{Password: npd_md5}
	err = tUser.Update(userID)
	if err != nil {
		return err
	}
	return nil
}

// 重置验证码到redis
func resetVerifyKey(userID uint, verifyCode string) error {
	redisKey := fmt.Sprintf(redisKeyFmt, strconv.Itoa(int(userID))) // 格式化存入的验证码
	oldKey, err := getVerifyKey(userID)
	if err == nil {
		_, err = redis.Do("del", redisKey) // 如果之前的key存在则先删除
		if err != nil {
			return fmt.Errorf("reset password key delete err: %s", err)
		}
	} else {
		if oldKey != "expired" {
			return err
		}
	}
	_, err = redis.Do("set", redisKey, verifyCode)
	if err != nil {
		return fmt.Errorf("reset password set err: %s", err)
	}
	_, err = redis.Do("expire", redisKey, duration)
	if err != nil {
		return fmt.Errorf("reset password set err: %s", err)
	}
	return nil
}

// 从redis读取验证码
func getVerifyKey(userID uint) (string, error) {
	redisKey := fmt.Sprintf(redisKeyFmt, strconv.Itoa(int(userID))) // 格式化存入的验证码
	res, err := rredis.String(redis.Do("get", redisKey))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			return "expired", fmt.Errorf("reset password is expired")
		}
		return "", fmt.Errorf("reset password get err: %s", err)
	}
	return res, nil
}
