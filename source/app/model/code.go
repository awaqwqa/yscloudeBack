package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义业务的状态码
type MyCode int64

const (
	CodeUnknowError     MyCode = 2401
	CodeSuccess         MyCode = 1000
	CodeInvalidParams   MyCode = 1001
	CodeUserExist       MyCode = 1002
	CodeUserNotExist    MyCode = 1003
	CodeInvalidPassword MyCode = 1004
	CodeServerBusy      MyCode = 1005

	CodeInvalidToken      MyCode = 1006
	CodeInvalidAuthFormat MyCode = 1007
	CodeNotLogin          MyCode = 1008
	CodeUserNameFalse     MyCode = 1009
	CodeUserPasswdFalse   MyCode = 1010
	CodeInvalidKey        MyCode = 1011
	CodeBindJsonFalse     MyCode = 1012
	CodeInvalidLevel      MyCode = 1013
)

var msgFlags = map[MyCode]string{
	CodeUnknowError:       "未知错误",
	CodeSuccess:           "success",
	CodeInvalidParams:     "请求参数错误",
	CodeUserExist:         "用户名重复",
	CodeUserNotExist:      "用户不存在",
	CodeInvalidPassword:   "用户名或者密码错误",
	CodeServerBusy:        "服务器繁忙",
	CodeInvalidToken:      "token无效",
	CodeInvalidAuthFormat: "认证格式有错误",
	CodeNotLogin:          "未登录",
	CodeUserNameFalse:     "名字非法",
	CodeUserPasswdFalse:   "密码非法",
	CodeInvalidKey:        "密钥不存在",
	CodeBindJsonFalse:     "绑定数据错误",
	CodeInvalidLevel:      "越权处理",
}

func BackError(ctx *gin.Context, code MyCode) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"code": code,
		"msg":  code.Msg(),
	})
}
func BackSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  CodeSuccess.Msg(),
	})
}
func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
