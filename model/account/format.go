// 结构体之间格式化
package account

import (
	"encoding/json"
	"fmt"
	"login-server/model/db/postgre"
	"time"
)

const time_format = "2006-01-02 15:04:05"

type UserInfoStr struct {
	Name     string
	Age      uint
	Birthday string
	Address  string
	Other    json.RawMessage
}

type UserStr struct {
	ID         uint
	LoginName  string
	Mobile     string
	Email      string
	Source     string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
	LoginTime  string
	LoginCount uint
}

// 从UserInfoStr解析到UserInfo，转化time和json类型
func ToUserInfo(userInfo *postgre.UserInfo, userInfoStr *UserInfoStr) error {
	var err error
	userInfo.Name = userInfoStr.Name
	userInfo.Age = userInfoStr.Age
	if userInfoStr.Birthday != "" {
		userInfo.Birthday, err = time.Parse(time_format, userInfoStr.Birthday)
		if err != nil {
			return fmt.Errorf("UserInfo.Birthday is not right format, format: %%Y-%%m-%%d %%H:%%M:%%S")
		}
	}
	userInfo.Address = userInfoStr.Address
	userInfo.Other = userInfoStr.Other
	return nil
}

// 从UserInfo解析到UserInfoStr
func ToUserInfoStr(userInfo *postgre.UserInfo, userInfoStr *UserInfoStr) error {
	userInfoStr.Name = userInfo.Name
	userInfoStr.Age = userInfo.Age
	userInfoStr.Birthday = userInfo.Birthday.Format(time_format)
	userInfoStr.Address = userInfo.Address
	userInfoStr.Other = userInfo.Other
	return nil
}

// 从User解析到UserStr
func ToUserStr(user *postgre.User, userStr *UserStr) error {
	userStr.ID = user.ID
	userStr.LoginName = user.LoginName
	userStr.Mobile = user.Mobile
	userStr.Email = user.Email
	userStr.Source = user.Source
	userStr.CreatedAt = user.CreatedAt.Format(time_format)
	userStr.UpdatedAt = user.UpdatedAt.Format(time_format)

	if user.DeletedAt != nil {
		userStr.DeletedAt = user.DeletedAt.Format(time_format)
	} else {
		userStr.DeletedAt = ""
	}
	userStr.LoginTime = user.LoginTime.Format(time_format)
	userStr.LoginCount = user.LoginCount
	return nil
}
