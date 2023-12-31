package route

import (
	"yscloudeBack/source/app/controller"
	"yscloudeBack/source/app/middleware"

	"github.com/gin-gonic/gin"
)

type handleFunction func(ctx *gin.Context)

/*
func getHandleAndRouter() map[string]handleFunction {
	return map[string]handleFunction{}
}
*/

func NewRegisterRoute(rg *gin.Engine) *RegisterRoute {
	return &RegisterRoute{
		RegisterEngine: rg,
	}
}

type RegisterRoute struct {
	RegisterEngine *gin.Engine
}

// Initialization is performed (执行) to connect the router to the handle function.
func (rg *RegisterRoute) Register(cm *controller.ControllerMannager) {
	baseGroup := rg.RegisterEngine.Group(BASE_PATH)
	{
		noticeGroup := baseGroup.Group("/notice")
		{
			noticeGroup.GET("/get_value", cm.GetNoticeValue())
		}
		keyGroup := baseGroup.Group("/key")
		{
			keyGroup.GET("/get_key_price", cm.GetKeyPrice())
		}
		// 登录相关
		logGroup := baseGroup.Group(LOGPATH)
		{
			// Register
			logGroup.POST(getRegisterUrl(), cm.Register())
			logGroup.POST("/login", cm.Login())
			logGroup.GET("/get_user_info", cm.GetUserInfo())
			//logGroup.POST(getLogoutUrl())
		}

		// 管理员权限相关
		adminGroup := baseGroup.Group("/admin")
		adminGroup.Use(middleware.JWTAuthMiddleware())
		adminGroup.Use(middleware.CheckAdmin())
		{
			adminGroup.GET("/get_user_name", cm.GetUserName())
			adminGroup.GET("/get_users", cm.GetUsers())
			adminGroup.POST("/register_key", cm.RegisterKey())
			adminGroup.GET("/get_keys", cm.GetKey())
			adminGroup.POST("/del_key", cm.DelKey())
			adminGroup.POST("/set_fbtoken", cm.SetFbToken())
			adminGroup.GET("/get_fbtokens", cm.GetFbTokens())
			adminGroup.POST("/del_fbtoken", cm.DelFbTokens())
			adminGroup.POST("/update_notice", cm.UpdateNotice())
			adminGroup.POST("/update_key_price", cm.UpdateKeyPrice())
			adminGroup.POST("/update_user_balance", cm.UpdateBalance())
		}
		// 文件相关
		StructGroup := baseGroup.Group(STRUCTPATH)
		StructGroup.Use(middleware.JWTAuthMiddleware())
		rg.RegisterEngine.MaxMultipartMemory = 8 << 20 // 8 MiB
		{
			//StructGroup.GET("/get_struts", cm.GetStructs())
			StructGroup.POST("/upload_strut", cm.UploadFile())
		}
		// 用户信息相关
		userGroup := baseGroup.Group("/user")
		userGroup.Use(middleware.JWTAuthMiddleware())
		{
			userGroup.GET("/get_keys", cm.GetUserKeys())
			userGroup.POST("/del_keys", cm.DelUserKey())
			userGroup.POST("/add_key", cm.AddUserKey())
			userGroup.POST("/buy_key", cm.BuyUserKey())
			userGroup.GET("/get_file_name", cm.GetUserFileName())
			userGroup.GET("/get_files", cm.GetFilesInfo())
			userGroup.GET("/get_file_groups", cm.GetFileGroups())
			userGroup.POST("/build_structure", cm.LoadHandler())
			userGroup.POST("/add_file_group", cm.AddFileGroup())
			userGroup.POST("/delete_file_group", cm.DeleteFileGroup())
			userGroup.POST("/delete_file", cm.DeleteFile())
			userGroup.GET("/get_balance", cm.GetBalance())
			userGroup.POST("/add_slot", cm.AddSlots())
			userGroup.POST("/update_slot", cm.UpdateSlots())
			userGroup.GET("/get_slots", cm.GetUserSlots())
			userGroup.POST("/del_slots", cm.DelUserSlots())
			userGroup.GET("/get_instance_status", cm.GetInstanceStatus())
			//userGroup.GET("/stream-instance-output", cm.GetStreamOutput())
		}
	}

}

// router wouldnt be imported .The router package is used to initialize the router similar a controller
func InitRoute(r *gin.Engine, cm *controller.ControllerMannager) {
	//跨域插件
	r.Use(Cors())
	rg := NewRegisterRoute(r)
	rg.Register(cm)
}
