// 权限管理handler
// 解析上报参数
// 调用github.com/login_server/authen中函数
package authen

import (
	"net/http"
	"strconv"
	"server/log"
	"github.com/gin-gonic/gin"
	"server/db/postgre"
	"server/base"
)

// 新增权限
func CreateHandler(c *gin.Context) {
	var right postgre.Right
	c.BindJSON(&right)
	err := AddRight(&right)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 删除权限
func DeleteHandler(c *gin.Context) {
	rightID, err := base.GetUriID("RightID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	err = DeleteRight(uint(rightID))
	if err == nil {
		base.RespondWithCode(204, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// RightID请求
type RightUpdate struct {
	Server string
	Name   string
	Desc   string
}

// 更新权限
func UpdateHandler(c *gin.Context) {
	rightID, err := base.GetUriID("RightID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	var rightUpdate RightUpdate
	c.BindJSON(&rightUpdate)
	log.Debug("RightUpdate: %v", rightUpdate)
	err = UpdateRight(uint(rightID), &postgre.Right{Server: rightUpdate.Server, Name: rightUpdate.Name, Desc: rightUpdate.Desc})
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 查询所有权限
func QueryAllHandler(c *gin.Context) {
	var rights *[]postgre.Right
	rights, err := QueryAll()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "", "data": rights})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": string(err.Error())})
	}
}

type ReqBindRight struct {
	UserID  uint
	RightID uint
}

// 绑定权限给用户
func BindRightHandler(c *gin.Context) {
	var reqBindRight ReqBindRight
	c.BindJSON(&reqBindRight)
	log.Debug("reqBindRight: %v\n", reqBindRight)
	err := BindRight(reqBindRight.UserID, reqBindRight.RightID)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 解除用户绑定的权限
func UnBindRightHandler(c *gin.Context) {
	rightID, err := base.GetUriID("RightID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	userID, err := base.GetUriID("UserID", c)
	if err != nil {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
	err = UnBindRight(userID, rightID)
	if err == nil {
		base.RespondWithCode(201, "", nil, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), nil, c)
	}
}

// 查看用户绑定的权限
func QueryBindRightHandler(c *gin.Context) {
	userIDStr := c.Param("UserID")
	log.Debug("UserID: %v\n", userIDStr)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		base.RespondWithCode(400, "UserID must be int", nil, c)
	}
	queryBindRight, err := QueryBindRight(uint(userID))
	data := map[string]interface{}{"data": queryBindRight}
	if err == nil {
		base.RespondWithCode(200, "", data, c)
	} else {
		base.RespondWithCode(400, string(err.Error()), data, c)
	}
}
