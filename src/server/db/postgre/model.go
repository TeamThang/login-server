// 表结构定义
// gorm
package postgre

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 用户基本表
type User struct {
	gorm.Model

	LoginName  string `gorm:"size:64;not null;index"` // 登录名, 未定义unique是为了支持软删除
	Password   string `gorm:"size:64;not null"`       // 密码
	Mobile     string `gorm:"size:20;index"`          // 手机号, 支持登录, 未定义unique是为了支持空值
	Email      string `gorm:"size:100;index"`         // 邮箱, 支持登录, 未定义unique是为了支持空值
	LoginTime  time.Time                              // 登录时间
	LoginCount uint   `gorm:"default:0"`              // 登录次数
	Source     string `gorm:"size:16"`                // 用户来源，注册所在服务

	UserInfo           UserInfo            `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"` // 用户详细信息表
	Emails             []Email             `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"` // 用户邮箱列表
	UserRightRelations []UserRightRelation `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"` // 关联用户权限表
	LoginLogs          []LoginLog          `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"` // 用户操作日志
	EmailLogs          []EmailLog          `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"` // 邮件发送日志
}

// 用户具体信息
type UserInfo struct {
	gorm.Model

	UserID   uint                     // 外键 (属于), tag `index`是为该列创建索引
	Name     string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
	Age      uint
	Birthday time.Time
	Address  string `gorm:"type:varchar(100)"`
	Other    []byte `gorm:"type:json;default:''"`
}

// 邮箱列表，用户其他邮箱
type Email struct {
	gorm.Model

	UserID     uint                                    // User.id外键
	Email      string `gorm:"type:varchar(100);index"` // 邮箱, 未定义unique是为了支持软删除
	Subscribed bool                                    // 是否订阅，订阅接受提醒邮件
}

// 验证表
type Verification struct {
	gorm.Model

	UserID  uint                             // User.id外键
	Content string `gorm:"size:64;not null"` // 验证的内容
	Status  bool                             // 验证的状态
}

// 用户权限关系表
type UserRightRelation struct {
	UserID  uint `gorm:"primary_key"` // 用户id
	RightID uint `gorm:"primary_key"` // 权限id
}

// 权限表
type Right struct {
	ID     uint   `gorm:"primary_key"`
	Server string `gorm:"size:16;not null"` // 支持服务名
	Name   string `gorm:"size:64;not null"` // 权限名称, 对应服务的权限分组
	Desc   string `gorm:"size:200"`         // 权限描述

	UserRightRelations []UserRightRelation `gorm:"ForeignKey:RightID;AssociationForeignKey:Refer"` // 关联用户权限表
}

// 登录日志表
type LoginLog struct {
	ID       uint   `gorm:"primary_key"`
	UserID   uint                                   // 操作人
	OpType   string `gorm:"size:20;index;not null"` // 操作类型: login / logout
	OpTime   time.Time                              // 操作时间
	Duration uint                                   // 登录持续时常
	Token    string `gorm:"size:64;index"`          // 登录token
}

// 邮件日志表
type EmailLog struct {
	ID          uint `gorm:"primary_key"`
	UserID      uint      // 操作人
	Receiver    string    // 收件人
	Subject     string    // 邮件主题
	ContentType string    // 邮件内容类型
	Content     string    // 邮件内容
	OpTime      time.Time // 操作时间
}
