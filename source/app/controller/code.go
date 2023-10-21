package controller

// 定义业务的状态码
type MyCode int64

const (
	CodeSuccess         MyCode = 1000
	CodeInvalidParams   MyCode = 1001
	CodeUserExist       MyCode = 1002
	CodeUserNotExist    MyCode = 1003
	CodeInvalidPassword MyCode = 1004
	CodeServerBusy      MyCode = 1005

	CodeInvalidToken      MyCode = 1006
	CodeInvalidAuthFormat MyCode = 1007
	CodeNotLogin          MyCode = 1008
)

var msgFlags = map[MyCode]string{
	CodeSuccess:           "success",
	CodeInvalidParams:     "请求参数错误",
	CodeUserExist:         "用户名重复",
	CodeUserNotExist:      "用户不存在",
	CodeInvalidPassword:   "用户名或者密码错误",
	CodeServerBusy:        "服务器繁忙",
	CodeInvalidToken:      "token无效",
	CodeInvalidAuthFormat: "认证格式有错误",
	CodeNotLogin:          "未登录",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}