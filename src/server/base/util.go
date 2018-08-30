package base

import (
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/log"
)

func RespondWithCode(code int, message string, otherData map[string]interface{}, c *gin.Context) {
	resp := map[string]interface{}{"status": code, "message": message}
	if otherData != nil {
		for k, v := range otherData {
			resp[k] = v
		}
	}
	c.JSON(code, resp)
	c.Abort()
}

// 从uri中获取id变量
func GetUriID(IDName string, c *gin.Context) (uint, error) {
	idStr := c.Param(IDName)
	log.Debug("%s: %v\n", IDName, idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("%s must be int", IDName)
	}
	return uint(id), nil
}

// 从uri中获取string变量
func GetUriStr(IDName string, c *gin.Context) (string, error) {
	res := c.Param(IDName)
	log.Debug("%s: %v\n", IDName, res)
	return res, nil
}
