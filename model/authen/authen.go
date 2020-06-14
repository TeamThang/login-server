// 权限配置和绑定用户
package authen

import (
	"fmt"
	"login-server/model/db/postgre"
	"login-server/util"
)

// 添加权限
func AddRight(right *postgre.Right) error {
	needParam := []string{"Server", "Name"}
	err := util.CheckAttrExist(right, needParam)
	if err != nil {
		return err
	}
	db := postgre.DB.Where(right).First(&postgre.Right{})
	if db.Error == nil {
		return fmt.Errorf("right is existed")
	}
	if !db.RecordNotFound() {
		return db.Error
	}
	db = postgre.DB.Create(right)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 删除权限
// 先删权限表，后删除用户权限关系表
func DeleteRight(rightID uint) error {
	if rightID == 0 {
		return fmt.Errorf("RightID is needed")
	}
	right := &postgre.Right{}
	db := postgre.DB.First(right, rightID)
	if db.Error != nil {
		if string(db.Error.Error()) != "record not found" {
			return fmt.Errorf("right is not existed")
		}
		return db.Error
	}
	db = postgre.DB.Unscoped().Delete(right) // 删除权限表
	if db.Error != nil {
		return db.Error
	}
	userRightRelation := &postgre.UserRightRelation{}
	db = postgre.DB.Where(&postgre.UserRightRelation{RightID: rightID}).First(userRightRelation)
	if db.Error != nil {
		if db.RecordNotFound() { // 没有记录无需删除
			return nil
		}
		return db.Error // 异常返回false和对应error
	}
	db = postgre.DB.Unscoped().Delete(userRightRelation) // 删除用户权限关系表
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 修改权限
func UpdateRight(rightID uint, tRight *postgre.Right) error {
	if rightID == 0 {
		return fmt.Errorf("RightID is needed")
	}
	right := &postgre.Right{}
	db := postgre.DB.First(right, rightID)
	if db.Error != nil {
		return db.Error
	}
	db = postgre.DB.Model(&right).Updates(tRight)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 查询所有权限
func QueryAll() (*[]postgre.Right, error) {
	rights := new([]postgre.Right)
	db := postgre.DB.Find(rights)
	if db.Error != nil {
		return nil, db.Error
	}
	return rights, nil
}

// 查询某个权限
func QueryOne(rightID uint) (*postgre.Right, error) {
	right := &postgre.Right{}
	db := postgre.DB.First(right, rightID)
	if db.Error != nil {
		return nil, db.Error
	}
	return right, nil
}

// 绑定权限到用户
func BindRight(userID uint, rightID uint) error {
	if userID == 0 {
		return fmt.Errorf("userID is needed")
	}
	if rightID == 0 {
		return fmt.Errorf("rightID is needed")
	}
	db := postgre.DB.First(&postgre.User{}, userID) // 检查userID
	if db.Error != nil {
		if string(db.Error.Error()) == "record not found" {
			return fmt.Errorf("userID is not existed")
		}
		return db.Error
	}
	db = postgre.DB.First(&postgre.Right{}, rightID) // 检查rightID
	if db.Error != nil {
		if db.RecordNotFound() {
			return fmt.Errorf("rightID is not existed")
		}
		return db.Error
	}
	userRightRelation := &postgre.UserRightRelation{UserID: userID, RightID: rightID}
	db = postgre.DB.Where(&postgre.UserRightRelation{UserID: userID, RightID: rightID}).First(&postgre.UserRightRelation{}) // 检查已有绑定
	if db.Error == nil {
		return fmt.Errorf("User[%v] and Right[%v] are binded", userID, rightID)
	}
	if string(db.Error.Error()) == "record not found" {
		if db = postgre.DB.Create(userRightRelation); db.Error == nil {
			return nil
		} else {
			return db.Error
		}
	}
	return db.Error
}

// 取消用户绑定的权限
func UnBindRight(userID uint, rightID uint) error {
	if userID == 0 {
		return fmt.Errorf("userID is needed")
	}
	if rightID == 0 {
		return fmt.Errorf("rightID is needed")
	}
	res := &postgre.UserRightRelation{}
	cond := &postgre.UserRightRelation{UserID: userID, RightID: rightID}
	db := postgre.DB.Where(cond).First(res)
	if db.Error != nil {
		return db.Error
	}
	db = postgre.DB.Unscoped().Delete(res) // 删除用户权限关系表
	if db.Error != nil {
		return db.Error
	}
	return nil
}

type RightQuery struct {
	UserID uint
	Rights []postgre.Right
}

// 查询用户绑定的权限
func QueryBindRight(userID uint) (*RightQuery, error) {
	if userID == 0 {
		return nil, fmt.Errorf("userID is needed")
	}
	rightQuery := &RightQuery{}
	rightQuery.UserID = userID
	res := []postgre.UserRightRelation{}
	db := postgre.DB.Where(&postgre.UserRightRelation{UserID: userID}).Find(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	for _, v := range res {
		right, err := QueryOne(v.RightID)
		if err != nil {
			return nil, err
		}
		rightQuery.Rights = append(rightQuery.Rights, *right)
	}
	return rightQuery, nil
}

// 根据Right.Server查询对应RightID
func QueryRightByServer(server string) (rightID uint, err error) {
	rightRes := &postgre.Right{}
	db := postgre.DB.Where(&postgre.Right{Server: server}).First(rightRes)
	if db.Error != nil {
		if db.RecordNotFound() {
			return 0, fmt.Errorf("right is not existed")
		}
		return 0, db.Error
	}
	return rightRes.ID, nil
}
