/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    utils
 * @Date:    2021/5/28 3:29 下午
 * @package: http
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package mygin

import (
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/errcode"
	"net/http"
)

func DecodeUrlVal(c *gin.Context, key string) (string, bool) {
	v, ok := c.GetQuery(key)
	if !ok {
		ErrInterrupt(c, errcode.InvalidParam)
	}
	return v, ok
}

func BindArgs(c *gin.Context, arg interface{}) bool {
	err := c.Bind(arg)
	if err != nil {
		ErrInterrupt(c, errcode.InvalidParam.WithErr(err))
		return false
	}
	return true
}

func BindQueryArgs(c *gin.Context, arg interface{}) bool {
	err := c.BindQuery(arg)
	if err != nil {
		ErrInterrupt(c, errcode.InvalidParam.WithErr(err))
		return false
	}
	return true
}

func BindJsonArgs(c *gin.Context, arg interface{}) bool {
	err := c.BindJSON(arg)
	if err != nil {
		ErrInterrupt(c, errcode.InvalidParam.WithErr(err))
		return false
	}
	return true
}

func BindHeader(c *gin.Context, arg interface{}) bool {
	err := c.BindHeader(arg)
	if err != nil {
		ErrInterrupt(c, errcode.InvalidParam.WithErr(err))
		return false
	}
	return true
}

func PkgMsgWrite(c *gin.Context, data interface{}) {
	code := errcode.Success
	dataMap := gin.H{"code": code.Code(), "msg": code.ErrMsg()}
	if data != nil {
		dataMap["data"] = data
	}
	c.JSON(http.StatusOK, dataMap)
}

func ErrInterrupt(c *gin.Context, err errcode.IErr) {
	c.JSON(http.StatusOK, gin.H{"code": err.Code(), "msg": err.ErrMsg()})
	c.Abort()
	c.Error(err)
}

func HasErr(c *gin.Context, errs ...error) bool {
	err := errcode.Errors(errs...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		c.Abort()
		c.Error(err)
		return true
	}
	return false
}
