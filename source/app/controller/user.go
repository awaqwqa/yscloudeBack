package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
	"yscloudeBack/source/app/middleware"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func checkName(name string) (bool, error) {
	// 检查长度
	if len(name) < 3 || len(name) > 16 {
		return false, fmt.Errorf("name length must be between 3 and 16")
	}

	// 检查是否符合命名规范（包含字母、中文字符和下划线）
	matched, err := regexp.MatchString("^[A-Za-z\\p{Han}_]+$", name)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, fmt.Errorf("name must only contain letters, Chinese characters, and underscores")
	}
	return true, nil
}

func checkPasswd(password string) error {
	if len(password) < 8 || len(password) > 20 {
		return fmt.Errorf("password must be between 8 and 20 characters")
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
func checkIsAllow(name string, passwd string, key string) (bool, model.MyCode) {
	if ok, _ := checkName(name); !ok {
		return false, model.CodeUserNameFalse
	}
	//检查是否合法
	if err := checkPasswd(passwd); err != nil {
		return false, model.CodeUserPasswdFalse
	}
	//检查语法是否合法
	if err := checkKey(key); err != nil {
		fmt.Println(err)
		return false, model.CodeInvalidKey
	}
	return true, 0
}
func (cm *ControllerMannager) GetUserName() gin.HandlerFunc {
	rg := cm.GetDbManager()
	return func(ctx *gin.Context) {
		users, err := rg.GetUsers()
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		resultSlice := []string{}
		for _, v := range users {
			resultSlice = append(resultSlice, v.UserName)
		}
		model.BackSuccess(ctx, resultSlice)
		return
	}
}
func (cm *ControllerMannager) GetUsers() gin.HandlerFunc {
	rg := cm.GetDbManager()
	return func(ctx *gin.Context) {
		users, err := rg.GetUsers()
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}

		model.BackSuccess(ctx, users)
		return
	}
}
func (cm *ControllerMannager) GetUserInfo() gin.HandlerFunc {
	db := cm.GetDbManager()
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		claim, err := utils.ParseToken(token)
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		userName := claim.UserName
		user, err := db.GetUserByUserName(userName)
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		roles := []string{"user"}
		if userName == "admin" {
			roles = append(roles, "admin")
		}
		model.BackSuccess(ctx, gin.H{
			"roles":        roles,
			"name":         userName,
			"balance":      user.Balance,
			"qq":           user.QQNumber,
			"avatar":       "nil",
			"introduction": "nil",
		})
		return

	}
}

// 获取userFileName
func (cm *ControllerMannager) GetUserFileName() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		name, err := middleware.GetContextName(ctx)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		user, err := manager.GetUserByUserName(name)
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		infoCopys, err := user.GetAllStructureInfoCopy()
		if err != nil {
			model.BackErrorByString(ctx, "cant get structures")
			return
		}
		model.BackSuccess(ctx, infoCopys)
		return
	}
}
func (cm *ControllerMannager) AddUserKey() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			Key       string `form:"key" binding:"required" json:"key"`
			FileGroup string `form:"file_group" binding:"required" json:"file_group"`
		}
		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		userKey, err := manager.GetKeyByValue(form.Key)
		if err != nil {
			model.BackError(ctx, model.CodeGetKeyFalse)
			return
		}
		// 检查userkey是否被使用
		if userKey.Status {
			model.BackError(ctx, model.CodeCodeIsUsed)
			return
		}
		// 判断userKey是否是期待类型
		if userKey.Usage != model.USAGE_LOAD {
			model.BackError(ctx, model.CodeCodeTypeFalse)
			return
		}
		// 获取 user
		name, isExit := ctx.Get(middleware.ContextName)
		if !isExit {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		value, isOk := name.(string)
		if !isOk {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		user, err := manager.GetUserByUserName(value)
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}

		err = manager.UpdateKeyStatus(userKey, true)
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		err = manager.UpdateKeyFileGroupName(userKey, form.FileGroup)
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		// 绑定user和key
		err = manager.AssociateKeyWithUser(user.ID, userKey.ID)
		if err != nil {
			utils.Error(err.Error())
			model.BackError(ctx, model.CodeUnknowError)
			return
		}

		model.BackSuccess(ctx, nil)
	}
}

func (cm *ControllerMannager) DelUserKey() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			DelKey string `form:"del_key" binding:"required" json:"del_key"`
		}

		if err := ctx.ShouldBind(&form); err != nil {
			model.BackError(ctx, model.CodeInvalidKey)
			return
		}
		userName, isok := ctx.Get(middleware.ContextName)
		if !isok {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		if value, isok := userName.(string); isok {
			user, err := manager.GetUserByUserName(value)
			if err != nil {
				model.BackError(ctx, model.CodeGetUserFalse)
				return
			}

			if !user.CheckLoadKey(form.DelKey) {
				model.BackError(ctx, model.CodeInvalidKey)
				return
			}
			err = manager.DeleteKey(form.DelKey)
			if err != nil {
				model.BackError(ctx, model.CodeUnknowError)
				return
			}
			model.BackSuccess(ctx, nil)
			return
		}
		model.BackSuccess(ctx, model.CodeUnknowError)
	}
}
func (cm *ControllerMannager) GetUserKeys() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		userName, isok := ctx.Get(middleware.ContextName)
		if !isok {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		value, isok := userName.(string)
		if !isok {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		user, err := manager.GetUserByUserName(value)
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}

		keys := user.GetLoadKeys()
		if keys != nil {
			model.BackSuccess(ctx, keys)
			return
		}
		model.BackSuccess(ctx, keys)

		return
	}
}

// Register 用户注册
func (cm *ControllerMannager) Register() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var rf *model.RegisterForm
		code, err := model.BindStruct(ctx, &rf)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		name := rf.UserName
		key := rf.RedeemKey
		passwd := rf.Password
		qq := rf.QQ
		if ok, code := checkIsAllow(name, passwd, key); !ok {
			model.BackError(ctx, code)
			return
		}
		//检查是否存在
		_, err = manager.GetUserByUserName(name)
		if err == nil {
			model.BackError(ctx, model.CodeUserExist)
			return
		}
		//检查key是否存在
		value, err := manager.GetKeyByValue(key)
		if err != nil {
			model.BackError(ctx, model.CodeInvalidKey)
			return
		}
		if value.Usage != model.USAGE_REGISTER {
			model.BackError(ctx, model.CodeCodeTypeFalse)
			return
		}
		//删除key
		err = manager.DeleteKey(key)
		if err != nil {
			model.BackError(ctx, model.CodeGetKeyFalse)
			return
		}
		// if pass  all
		user := model.NewUser(name, utils.Md5Encrypt(passwd), key)
		//获取token
		token, _, err := utils.GenToken(name)
		if err != nil {
			model.BackError(ctx, model.CodeGetTokenFalse)
			return
		}
		user.Token = token
		user.QQNumber = qq
		//存入数据库
		err = manager.CreateUser(user)
		if err != nil {
			model.BackError(ctx, model.CodeCreateUserFalse)
			return
		}
		//返回成功信息
		ctx.JSON(http.StatusOK, gin.H{
			"code":  model.CodeSuccess,
			"msg":   model.CodeSuccess.Msg(),
			"Token": token,
		})
		return
	}
}
func (cm *ControllerMannager) Login() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var lf *model.LoginForm
		if err := ctx.ShouldBindJSON(&lf); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if ok {
				model.BackError(ctx, model.CodeUnknowError)
				return
			}
		}
		if ok, _ := checkName(lf.UserName); !ok {
			model.BackError(ctx, model.CodeUserNameFalse)
			return
		}
		if err := checkPasswd(lf.Password); err != nil {
			model.BackError(ctx, model.CodeUserPasswdFalse)
			return
		}
		user, err := manager.GetUserByUserName(lf.UserName)
		if err != nil {
			model.BackError(ctx, model.CodeUserNotExist)
			return
		}
		if user.Password != utils.Md5Encrypt(lf.Password) {
			model.BackError(ctx, model.CodeInvalidPassword)
			return
		}
		token, _, err := utils.GenToken(lf.UserName)
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":  model.CodeSuccess,
			"msg":   model.CodeSuccess.Msg(),
			"token": token,
		})

	}
}
