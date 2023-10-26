package controller

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func checkName(name string) (bool, error) {
	// 检查长度
	if len(name) < 3 || len(name) > 8 {
		return false, fmt.Errorf("name length must be between 3 and 8")
	}

	// 检查是否符合命名规范（这里假设名字只能包含字母）
	matched, err := regexp.MatchString("^[A-Za-z]+$", name)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, fmt.Errorf("name must only contain letters")
	}
	return true, nil
}

func checkPasswd(password string) error {
	if len(password) < 8 || len(password) > 16 {
		return fmt.Errorf("password must be between 8 and 16 characters")
	}
	// 检查密码是否包含至少一个数字
	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return fmt.Errorf("password must include at least one digit")
	}

	// 检查密码是否包含至少一个字母
	if matched, _ := regexp.MatchString(`[A-Za-z]`, password); !matched {
		return fmt.Errorf("password must include at least one letter")
	}

	// 检查密码是否包含至少一个特殊字符
	if matched, _ := regexp.MatchString(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",<>\./?\\|`+"`"+`]`, password); !matched {
		return fmt.Errorf("password must include at least one special character")
	}

	return nil
}
func checkKey(key string) error {
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 characters long")
	}

	// 使用正则表达式来检查字符串是否只包含数字和字母
	matched, err := regexp.MatchString(`^[A-Za-z0-9]+$`, key)
	if err != nil {
		return fmt.Errorf("failed to check the key: %v", err)
	}
	if !matched {
		return fmt.Errorf("key must consist of letters and numbers only")
	}

	return nil
}
func checkIsAllow(name string, passwd string, key string) (bool, MyCode) {
	if ok, _ := checkName(name); !ok {
		return false, CodeUserNameFalse
	}
	//检查是否合法
	if err := checkPasswd(passwd); err != nil {
		return false, CodeUserPasswdFalse
	}
	//检查语法是否合法
	if err := checkKey(key); err != nil {
		return false, CodeInvalidKey
	}
	return true, 0
}

// Register 用户注册
func Register(manager *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rf *model.RegisterForm
		if err := ctx.ShouldBindJSON(&rf); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if ok {
				BackError(ctx, CodeUnknowError)
				return
			}
		}
		name := rf.UserName
		key := rf.RedeemKey
		passwd := rf.Password
		qq := rf.QQ
		if ok, code := checkIsAllow(name, passwd, key); !ok {
			BackError(ctx, code)
		}
		//检查是否存在
		_, err := manager.GetUserByUserName(name)
		if err == nil {
			BackError(ctx, CodeUserExist)
			return
		}
		//检查key是否存在
		if ok, _ := manager.CheckKey(key); !ok {
			BackError(ctx, CodeInvalidKey)
			return
		}

		// if pass  all
		user := model.NewUser(name, utils.Md5Encrypt(passwd), key)
		//获取token
		token, _, err := utils.GenToken(name)
		if err != nil {
			return
		}
		user.Token = token
		user.QQNumber = qq
		//存入数据库
		err = manager.CreateUser(user)
		if err != nil {
			return
		}
		//返回成功信息
		ctx.JSON(http.StatusOK, gin.H{
			"code":  CodeSuccess,
			"msg":   CodeSuccess.Msg(),
			"Token": token,
		})
		return
	}
}
func Login(manager *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lf *model.LoginForm
		if err := ctx.ShouldBindJSON(&lf); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if ok {
				BackError(ctx, CodeUnknowError)
				return
			}
		}
		if ok, _ := checkName(lf.UserName); !ok {
			BackError(ctx, CodeUserNameFalse)
			return
		}
		if err := checkPasswd(lf.UserName); err != nil {
			BackError(ctx, CodeUserPasswdFalse)
			return
		}
		user, err := manager.GetUserByUserName(lf.UserName)
		if err != nil {
			BackError(ctx, CodeUserNotExist)
			return
		}
		if user.Password != utils.Md5Encrypt(lf.Password) {
			BackError(ctx, CodeInvalidPassword)
			return
		}
		token, _, err := utils.GenToken(lf.UserName)
		if err != nil {
			BackError(ctx, CodeUnknowError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":  CodeSuccess,
			"msg":   CodeSuccess.Msg(),
			"Token": token,
		})

	}
}
