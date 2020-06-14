// 封装特殊和需多次调用的数据库操作
// gorm
// 用model中定义的结构体传递数据
package postgre

import (
	"fmt"
	"reflect"
	"sync"
)

// 条件查询数据库，返回是否存在
// 基于gorm,condWhere: where条件; condNot: not条件
func QueryExist(condWhere interface{}, condNot interface{}) (bool, error) {
	if condWhere == nil && condNot == nil {
		return false, nil
	}
	outType := reflect.ValueOf(condWhere).Type().Elem()
	out := reflect.New(outType).Elem().Addr().Interface() // 创建和cod类型相同的输出变量
	db := DB.Where(condWhere).Not(condNot).First(out)
	if db.Error != nil {
		if db.RecordNotFound() { // 没有记录返回false
			return false, nil
		}
		return false, db.Error // 异常返回false和对应error
	}
	return true, nil
}

// 启动多个goroutine检查数据是否存在
func queryExists(condWhere interface{}, condNot interface{}, ch chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := QueryExist(condWhere, condNot)
	if err != nil {
		panic(err)
	}
	ch <- res
	close(ch)
}

// 并发校验LoginName,Mobile,Email是否存在
func checkBasicElems(user *User) error {
	chLoginName := make(chan bool, 1)
	chMobile := make(chan bool, 1)
	chEmail := make(chan bool, 1)
	var wg sync.WaitGroup
	if user.LoginName != "" { // 空值不算重复
		wg.Add(1)
		go queryExists(&User{LoginName: user.LoginName}, &User{}, chLoginName, &wg)
	} else {
		chLoginName <- false
		close(chLoginName)
	}
	if user.Mobile != "" {
		wg.Add(1)
		go queryExists(&User{Mobile: user.Mobile}, &User{}, chMobile, &wg)
	} else {
		chMobile <- false
		close(chMobile)
	}
	if user.Email != "" {
		wg.Add(1)
		go queryExists(&User{Email: user.Email}, &User{}, chEmail, &wg)
	} else {
		chEmail <- false
		close(chEmail)
	}
	wg.Wait()
	if <-chLoginName {
		return fmt.Errorf("LoginName(%q) is already existed", user.LoginName)
	}
	if <-chMobile {
		return fmt.Errorf("Mobile(%q) is already existed", user.Mobile)
	}
	if <-chEmail {
		return fmt.Errorf("Email(%q) is already existed", user.Email)
	}
	return nil
}

// 创建用户User表和UserInfo表
// 事务处理保证两个表创建成功
func UserCreate(user *User, userInfo *UserInfo) (userID uint, err error) {
	err = checkBasicElems(user)
	if err != nil {
		return 0, err
	}
	db := DB
	tx := db.Begin()
	res := tx.Create(&user)
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}
	newUser, ok := res.Value.(**User) // 用事务后，返回类型变成地址的地址了～
	if !ok {
		tx.Rollback()
		return 0, fmt.Errorf("create user return err")
	}
	userInfo.UserID = (*newUser).ID
	if userInfo != nil {
		if err := tx.Create(&userInfo).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return (*newUser).ID, nil
}

// 修改User
// userID: User表主键id
func (tUser *User) Update(userID uint) error {
	var user User
	db := DB.First(&user, userID)
	if db.Error != nil {
		return db.Error
	}
	if err := checkBasicElems(tUser); err != nil {
		return err
	}
	db = DB.Model(&user).Updates(tUser)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 删除User
// userID: User表主键id
func UserDelete(userID uint) error {
	user := &User{}
	db := DB.First(&user, userID)
	if db.Error != nil {
		return db.Error
	}
	db = DB.Delete(user)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 查询所有User
func (user *User) QueryUser() (*[]User, error) {
	var users *[]User
	db := DB.Find(users)
	if db.Error != nil {
		return nil, db.Error
	}
	return users, nil
}

// 插入UserInfo
func (userInfo *UserInfo) Create(userID uint) error {
	userInfo.UserID = userID
	db := DB.Create(userInfo)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 修改UserInfo
func (tUserInfo *UserInfo) Update(userID uint) error {
	var userInfo UserInfo
	db := DB.Where(&UserInfo{UserID: userID}).First(&userInfo)
	if db.Error != nil {
		return db.Error
	}
	db = DB.Model(&userInfo).Updates(tUserInfo)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 删除UserInfo表
func UserInfoDelete(userID uint) error {
	userInfo := &UserInfo{}
	db := DB.Where(&UserInfo{UserID: userID}).First(userInfo)
	if db.Error != nil {
		if db.RecordNotFound() {
			return nil // 不存在就不需要删除
		}
		return db.Error
	}
	db = DB.Delete(userInfo)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// 删除UserRightRelations表
func UserRightRelationsDelete(userID uint) error {
	userRightRelations := []UserRightRelation{}
	db := DB.Where(&UserRightRelation{UserID: userID}).Find(&userRightRelations)
	if db.Error != nil {
		if db.RecordNotFound() {
			return nil // 不存在就不需要删除
		}
		return db.Error
	}
	for _, userRightRelation := range userRightRelations {
		db = DB.Unscoped().Delete(&userRightRelation)
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

// 删除Email表
func EmailDelete(userID uint) error {
	emails := []Email{}
	db := DB.Where(&Email{UserID: userID}).Find(&emails)
	if db.Error != nil {
		if db.RecordNotFound() {
			return nil // 不存在就不需要删除
		}
		return db.Error
	}
	for _, email := range emails {
		db = DB.Delete(&email)
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

// 删除Verify表
func VerifyDelete(userID uint) error {
	verfiys := []Verification{}
	db := DB.Where(&Email{UserID: userID}).Find(&verfiys)
	if db.Error != nil {
		if db.RecordNotFound() {
			return nil // 不存在就不需要删除
		}
		return db.Error
	}
	for _, verify := range verfiys {
		db = DB.Delete(&verify)
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

// 通过主键id查询UserInfo
// userID: User表主键id
func (userInfo *UserInfo) QueryUserInfoByUserID(userID uint) (*UserInfo, error) {
	db := DB.Where(&UserInfo{UserID: userID}).First(userInfo)
	if db.Error != nil {
		if db.RecordNotFound() {
			return nil, nil
		}
		return nil, db.Error
	}
	return userInfo, nil
}

// 通过userID查询用户表
func GetUserByID(userID uint) (*User, error) {
	var user User
	db := DB.First(&user, userID)
	if db.Error != nil {
		if db.RecordNotFound() { // 没有记录返回false
			return nil, fmt.Errorf("user[%d] is not existed", userID)
		}
		return nil, fmt.Errorf("get user failed, err: %v", db.Error)
	}
	return &user, nil
}

// 检查用户ID是否存在
// 返回nil表示存在
func CheckUserExist(userID uint) (*User, error) {
	user := &User{}
	db := DB.First(user, userID)
	if db.Error != nil {
		if string(db.Error.Error()) == "record not found" {
			return nil, fmt.Errorf("user is not existed")
		}
		return nil, db.Error
	}
	return user, nil
}

// 通过主键id查询User
// userID: User表主键id
func (user *User) QueryUserById(userID uint) (*User, error) {
	db := DB.First(user, userID)
	if db.Error != nil {
		return nil, db.Error
	}
	return user, nil
}
