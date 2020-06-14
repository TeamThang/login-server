// 登录权限验证
package login

import (
	"fmt"
	"login-server/model/authen"
	"login-server/model/db/postgre"
)

// 验证权限
func checkRight(userID uint, right *Right) error {
	if right.Server == "" {
		return fmt.Errorf("Login Right.Server is needed")
	}
	grantRights, err := authen.QueryBindRight(userID)
	if err != nil {
		return err
	}
	err = checkRightServer(&grantRights.Rights, right)
	if err != nil {
		return err
	}
	err = checkRightName(&grantRights.Rights, right)
	if err != nil {
		return err
	}
	return nil
}

// 检查支持的Server
// 配置all支持所有
func checkRightServer(grantRights *[]postgre.Right, reqRight *Right) error {
	flag := true
	for _, v := range *grantRights {
		if v.Server == "all" { // 包含all则支持所有服务登录
			flag = false
			break
		}
		if reqRight.Server == v.Server {
			flag = false
		}
	}
	if flag {
		return fmt.Errorf("user has no right to login server[%s]", reqRight.Server)
	}
	return nil
}

// 检查对应服务的权限
// 配置all支持所有
func checkRightName(grantRights *[]postgre.Right, reqRight *Right) error {
	flag := true
	for _, v := range *grantRights {
		if v.Name == "all" { // 包含all则支持所有服务登录
			flag = false
			break
		}
		if reqRight.Name == v.Name {
			flag = false
		}
	}
	if flag {
		return fmt.Errorf("user has no right to login server[%s] with Right[%s]", reqRight.Server, reqRight.Name)
	}
	return nil
}

// 坚持支持login_server权限
// admin(all和login)
func checkLoginRight(userID uint) (uint, error) {
	grantRights, err := authen.QueryBindRight(userID)
	if err != nil {
		return 0, err
	}
	flag := false
	for _, r := range grantRights.Rights {
		if r.Server == "all" { // Right.Server为all
			flag = true
			break
		}
		if r.Server == "login" { // Right.Server为login
			flag = true
			break
		}
	}
	if flag {
		return 1, nil
	}
	return 0, nil
}
