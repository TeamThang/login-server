// 邮箱验证
package verify

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
	rredis "github.com/gomodule/redigo/redis"

	"server/db/postgre"
	"server/db/redis"
	"server/conf"
	"server/util"
)

const (
	verifyKeyFmt string = "VERIFY_EMAIL:%s"
	duration            = 300
)

// 发送邮件验证码
func SendVerifyEmail(userID uint) error {
	if userID == 0 {
		return fmt.Errorf("UserID is needed")
	}
	user, err := postgre.GetUserByID(userID)
	if err != nil {
		return err
	}
	verifyCode := util.RandStringRunes(5)
	err = resetVerifyKey(user.ID, verifyCode)
	if err != nil {
		return err
	}

	config := conf.Config.Email
	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Email Verify")
	m.SetBody("text/html", makeEmail("text/html", user.LoginName, verifyCode))

	d := gomail.NewDialer(config.Host, config.Port, config.User, config.PassWord)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send email err: %v", err)
	}
	return nil
}

// 检查验证码
func CheckVerifyCode(userID uint, verifyCode string) (bool, error) {
	if userID == 0 {
		return false, fmt.Errorf("UserID is needed")
	}
	if verifyCode == "" {
		return false, fmt.Errorf("VerifyCode is needed")
	}
	vc, err := getVerifyKey(userID)
	if err != nil {
		return false, err
	}
	if vc == verifyCode {
		go func() {
			postgre.DB.Create(&postgre.Verification{UserID: userID, Content: "email", Status: true})
		}()
		return true, nil
	}
	return false, nil
}

func makeEmail(contentType string, loginName string, verifyCode string) string {
	var res string
	switch contentType {
	case "text/html":
		res = fmt.Sprintf(TempEmailHtml, loginName, loginName, verifyCode)
	}
	return res
}

// 重置验证码到redis
func resetVerifyKey(userID uint, verifyCode string) error {
	verifyKey := fmt.Sprintf(verifyKeyFmt, strconv.Itoa(int(userID))) // 格式化存入的验证码 VERIFY_EMAIL:userID
	oldKey, err := getVerifyKey(userID)
	if err == nil {
		_, err = redis.Do("del", verifyKey) // 如果之前的key存在则先删除
		if err != nil {
			return fmt.Errorf("verify key delete err: %s", err)
		}
	} else {
		if oldKey != "expired" {
			return err
		}
	}
	_, err = redis.Do("set", verifyKey, verifyCode)
	if err != nil {
		return fmt.Errorf("verify key set err: %s", err)
	}
	_, err = redis.Do("expire", verifyKey, duration)
	if err != nil {
		return fmt.Errorf("verify key set err: %s", err)
	}
	return nil
}

// 从redis读取验证码
func getVerifyKey(userID uint) (string, error) {
	verifyKey := fmt.Sprintf(verifyKeyFmt, strconv.Itoa(int(userID))) // 格式化存入的验证码 VERIFY_EMAIL:userID
	res, err := rredis.String(redis.Do("get", verifyKey))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			return "expired", fmt.Errorf("verify key is expired")
		}
		return "", fmt.Errorf("verify key get err: %s", err)
	}
	return res, nil
}
