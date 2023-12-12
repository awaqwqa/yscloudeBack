package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yscloudeBack/source/app/utils"
)

// 定义业务的状态码
type MyCode int64

const (
	CodeUnknowError     MyCode = 2401
	CodeError           MyCode = 2000
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
	CodeGetKeyFalse       MyCode = 1014
	CodeGetTokenFalse     MyCode = 1015
	CodeCreateUserFalse   MyCode = 1016
	CodeGetUserFalse      MyCode = 1017
	CodeCreateKeyFalse    MyCode = 1018
	CodeUpdateUserFalse   MyCode = 1019
	CodeCodeTypeFalse     MyCode = 1020
	CodeBindFalse         MyCode = 1021
	CodeCodeIsUsed        MyCode = 1022
	CodeGetFbTokenFalse   MyCode = 1023
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
	CodeInvalidKey:        "密钥无效",
	CodeBindJsonFalse:     "绑定数据错误",
	CodeInvalidLevel:      "越权处理",
	CodeGetKeyFalse:       "获取key失败",
	CodeGetTokenFalse:     "获取token失败",
	CodeCreateUserFalse:   "创建user失败",
	CodeGetUserFalse:      "获取user失败",
	CodeCreateKeyFalse:    "创建密钥失败",
	CodeUpdateUserFalse:   "更新user失败",
	CodeCodeTypeFalse:     "密钥类型错误",
	CodeBindFalse:         "server解析失败",
	CodeCodeIsUsed:        "密钥已被注册",
	CodeGetFbTokenFalse:   "fbToken未设置",
}

func BackError(ctx *gin.Context, code MyCode) {
	utils.Error(code.Msg())
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  code.Msg(),
	})
}
func BackErrorByString(ctx *gin.Context, titleString string) {
	utils.Error(titleString)
	ctx.JSON(http.StatusOK, gin.H{
		"code": CodeError,
		"msg":  titleString,
	})
}
func BackSuccess(ctx *gin.Context, body interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  CodeSuccess.Msg(),
		"body": body,
	})
}
func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
