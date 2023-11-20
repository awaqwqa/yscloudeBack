package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

const (
	ContextName = "username"
)

// JWTAuthMiddleware 基于JWT的认证中间件 负责处理操作
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": model.CodeInvalidToken,
				"msg":  model.CodeInvalidToken.Msg(),
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(ContextName, mc.UserName)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
func CheckAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, ok := ctx.Get("username")
		if !ok {
			model.BackError(ctx, model.CodeUnknowError)
			ctx.Abort()
			return
		}
		userName := value.(string)
		if userName != "admin" {
			model.BackError(ctx, model.CodeInvalidLevel)
			ctx.Abort()
			return
		}

	}
}
func GetContextName(ctx *gin.Context) (string, error) {
	value, isFind := ctx.Get(ContextName)
	if !isFind {
		return "", fmt.Errorf("cant find userName")
	}
	switch value.(type) {
	case string:
		name := value.(string)
		return name, nil
	default:
		return "", fmt.Errorf("cant find userName")
	}
}
