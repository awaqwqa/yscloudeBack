package route

import (
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/controller"
	"yscloudeBack/source/app/middleware"
	"yscloudeBack/source/app/model"
)

type handleFunction func(ctx *gin.Context)

/*
func getHandleAndRouter() map[string]handleFunction {
	return map[string]handleFunction{}
}
*/

func NewRegisterRoute(rg *gin.Engine, manager *model.DbManager) *RegisterRoute {
	return &RegisterRoute{
		RegisterEngine: rg,
		Db:             manager,
	}
}

type RegisterRoute struct {
	RegisterEngine *gin.Engine
	Db             *model.DbManager
}

// Initialization is performed (执行) to connect the router to the handle function.
func (rg *RegisterRoute) Register() {
	// TODO : connect the routers based on base_path to the handlefunction
	/*err := rg.RegisterMiddleware()
	if err != nil {
		return
	}*/
	baseGroup := rg.RegisterEngine.Group(BASE_PATH)
	{
		logGroup := baseGroup.Group(LOGPATH)
		{
			// Register
			logGroup.POST(getRegisterUrl(), controller.Register(rg.Db))
			//logGroup.POST(getLoginUrl(), controller.Login)
			//logGroup.POST(getLogoutUrl())
		}

		// key controller
		keyGroup := baseGroup.Group(KEYPATH)
		keyGroup.Use(middleware.JWTAuthMiddleware())
		keyGroup.Use(middleware.CheckAdmin())
		{
			keyGroup.POST("/register", controller.RegisterKey(rg.Db))
		}

		LoadGroup := baseGroup.Group(LOADPATH)
		LoadGroup.Use(middleware.JWTAuthMiddleware())
		{
			LoadGroup.POST(LOADSTAR, controller.LoadHandler(rg.Db))
		}
		StructGroup := baseGroup.Group(STRUCTPATH)
		StructGroup.Use(middleware.JWTAuthMiddleware())
		rg.RegisterEngine.MaxMultipartMemory = 8 << 20 // 8 MiB
		{
			StructGroup.GET(GETSTRUCT, controller.GetStruct(rg.Db))
			StructGroup.POST(UPLOADPATH, controller.UploadFile(rg.Db))
		}

	}

}

// router wouldnt be imported .The router package is used to initialize the router similar a controller
func InitRoute(r *gin.Engine, manager *model.DbManager) {
	//跨域插件
	r.Use(Cors())
	rg := NewRegisterRoute(r, manager)
	rg.Register()
}
