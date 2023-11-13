package route

import (
	"yscloudeBack/source/app/controller"
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/middleware"

	"github.com/gin-gonic/gin"
)

type handleFunction func(ctx *gin.Context)

/*
func getHandleAndRouter() map[string]handleFunction {
	return map[string]handleFunction{}
}
*/

func NewRegisterRoute(rg *gin.Engine, manager *db.DbManager) *RegisterRoute {
	return &RegisterRoute{
		RegisterEngine: rg,
		Db:             manager,
	}
}

type RegisterRoute struct {
	RegisterEngine *gin.Engine
	Db             *db.DbManager
}

// Initialization is performed (执行) to connect the router to the handle function.
func (rg *RegisterRoute) Register() {
	baseGroup := rg.RegisterEngine.Group(BASE_PATH)
	{
		// 登录相关
		logGroup := baseGroup.Group(LOGPATH)
		{
			// Register
			logGroup.POST(getRegisterUrl(), controller.Register(rg.Db))
			logGroup.POST("/login", controller.Login(rg.Db))
			logGroup.GET("/get_user_info", controller.GetUserInfo(rg.Db))
			//logGroup.POST(getLogoutUrl())
		}
		// 管理员权限相关
		adminGroup := baseGroup.Group("/admin")
		adminGroup.Use(middleware.JWTAuthMiddleware())
		adminGroup.Use(middleware.CheckAdmin())
		{
			adminGroup.GET("/get_user_name", controller.GetUserName(rg.Db))
			adminGroup.GET("/get_users", controller.GetUsers(rg.Db))
			adminGroup.POST("/register_key", controller.RegisterKey(rg.Db))
			adminGroup.GET("/get_keys", controller.GetKey(rg.Db))
			adminGroup.POST("/del_key", controller.DelKey(rg.Db))
		}
		// 导入相关
		LoadGroup := baseGroup.Group(LOADPATH)
		LoadGroup.Use(middleware.JWTAuthMiddleware())
		{
			LoadGroup.POST(LOADSTAR, controller.LoadHandler(rg.Db))
		}
		// 文件相关
		StructGroup := baseGroup.Group(STRUCTPATH)
		StructGroup.Use(middleware.JWTAuthMiddleware())
		rg.RegisterEngine.MaxMultipartMemory = 8 << 20 // 8 MiB
		{
			StructGroup.GET(GETSTRUCT, controller.GetStruct(rg.Db))
			StructGroup.POST(UPLOADPATH, controller.UploadFile(rg.Db))
		}
		// 用户信息相关
		userGroup := baseGroup.Group("/user")
		userGroup.Use(middleware.JWTAuthMiddleware())
		{
			userGroup.GET("/get_keys", controller.GetUserKeys(rg.Db))
			userGroup.POST("/del_keys", controller.DelUserKey(rg.Db))
			userGroup.POST("/add_key", controller.AddUserKey(rg.Db))
			userGroup.GET("/get_file_name", controller.GetUserFileName(rg.Db))
			userGroup.POST("/build_struct", controller.LoadHandler(rg.Db))
		}
	}

}

// router wouldnt be imported .The router package is used to initialize the router similar a controller
func InitRoute(r *gin.Engine, manager *db.DbManager) {
	//跨域插件
	r.Use(Cors())
	rg := NewRegisterRoute(r, manager)
	rg.Register()
}
