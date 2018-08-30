// 邮件提醒
package notify

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
	"server/db/postgre"
	"server/util"
	"server/conf"
	"server/log"
)

const time_format = "2006-01-02 15:04:05"

// 添加提醒邮箱
func AddEmail(userID uint, email string, sub bool) error {
	if userID == 0 {
		return fmt.Errorf("UserID is needed")
	}
	if email == "" {
		return fmt.Errorf("Email is needed")
	}
	_, err := postgre.CheckUserExist(userID)
	if err != nil {
		return err
	}
	err = util.CheckEmailFormat(email)
	if err != nil {
		return fmt.Errorf("please use right email address")
	}
	db := postgre.DB.Where(&postgre.Email{UserID: userID, Email: email}).First(&postgre.Email{})
	if db.Error == nil {
		return fmt.Errorf("email is existed")
	}
	if db.RecordNotFound()  { // 不存在在添加
		db = postgre.DB.Create(&postgre.Email{UserID: userID, Email: email, Subscribed: sub})
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

// 删除提醒邮箱
func DelEmail(emailID uint) error {
	if emailID == 0 {
		return fmt.Errorf("EmailID is needed")
	}
	email, err := checkEmailExist(emailID)
	if err != nil {
		return err
	}
	db := postgre.DB.Delete(email)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 设置邮箱为订阅
func SubEmail(emailID uint) error {
	if emailID == 0 {
		return fmt.Errorf("EmailID is needed")
	}
	email, err := checkEmailExist(emailID)
	if err != nil {
		return err
	}
	db := postgre.DB.Model(email).Update("Subscribed", true)
	if db.Error != nil {
		return fmt.Errorf("sub email err: %v", db.Error)
	}
	return nil
}

// 设置邮箱为非订阅
func UnSubEmail(emailID uint) error {
	if emailID == 0 {
		return fmt.Errorf("EmailID is needed")
	}
	email, err := checkEmailExist(emailID)
	if err != nil {
		return err
	}
	db := postgre.DB.Model(email).Update("Subscribed", false)
	if db.Error != nil {
		return fmt.Errorf("unsub email err: %v", db.Error)
	}
	return nil
}

type EmaiQuery struct {
	ID         uint
	CreatedAt  string
	UpdatedAt  string
	UserID     uint
	Email      string
	Subscribed bool
}


// 获得所有提醒邮箱
func QueryAllEmail(userID uint) (*[]EmaiQuery, error) {
	var emails []postgre.Email
	db := postgre.DB.Where(&postgre.Email{UserID: userID}).Find(&emails)
	if db.Error != nil {
		return nil, fmt.Errorf("query all notify email err: %v", db.Error)
	}
	res := make([]EmaiQuery, 0)
	for _, email := range emails {
		d := EmaiQuery{
			ID: email.ID,
			CreatedAt: email.CreatedAt.Format(time_format),
			UpdatedAt: email.UpdatedAt.Format(time_format),
			UserID: email.UserID,
			Email: email.Email,
			Subscribed: email.Subscribed,
		}
		res = append(res, d)

	}
	return &res, nil
}

// 获得需要提醒的邮箱
func GetNotifyEmail(userID uint) (*[] string, error) {
	var ret []string
	emails, err := QueryAllEmail(userID)
	if err != nil {
		return nil, err
	}
	for _, email := range *emails {
		if email.Subscribed {
			ret = append(ret, email.Email)
		}
	}
	return &ret, nil
}

// 发送提醒邮件
func SendNotifyEmail(userID uint, subject string, contentType string, content string) error {
	emails, err := GetNotifyEmail(userID)
	if err != nil {
		return err
	}
	if len(*emails) == 0 {
		return fmt.Errorf("notify email is not added")
	}
	err = SendEmail(userID, subject, contentType, content, emails)
	if err != nil {
		return err
	}
	return nil
}

// 发送邮件
func SendEmail(userID uint, subject string, contentType string, content string, emails *[] string) error {
	if userID == 0 {
		return fmt.Errorf("UserID is needed")
	}
	if subject == "" {
		return fmt.Errorf("Subject is needed")
	}
	if contentType == "" {
		return fmt.Errorf("ContentType is needed")
	}
	if content == "" {
		return fmt.Errorf("Content is needed")
	}
	config := conf.Config.Email
	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", *emails...)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, content)

	d := gomail.NewDialer(config.Host, config.Port, config.User, config.PassWord)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send email err: %v", err)
	}
	go func() {
		emailLog := &postgre.EmailLog{
			UserID:      userID,
			Receiver:    strings.Join(*emails, ";"),
			Subject:     subject,
			ContentType: contentType,
			Content:     content,
			OpTime:      time.Now(),
		}
		db := postgre.DB.Create(emailLog)
		if db.Error != nil {
			log.Error("write email log err: %v", db.Error)
		}
	}()
	return nil
}

// 检查邮箱ID是否存在
// 返回nil表示存在
func checkEmailExist(emailID uint) (*postgre.Email, error) {
	email := &postgre.Email{}
	db := postgre.DB.First(email, emailID)
	if db.Error != nil {
		if string(db.Error.Error()) == "record not found" {
			return nil, fmt.Errorf("email is not existed")
		}
		return nil, db.Error
	}
	return email, nil
}
