package route

import (
	"github.com/gin-gonic/gin"
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
	err := rg.RegisterLogRoute()
	if err != nil {
		return
	}
	err = rg.RegisterLoadRoute()
	if err != nil {
		return
	}

}

// router wouldnt be imported .The router package is used to initialize the router similar a controller
func InitRoute(r *gin.Engine, manager *model.DbManager) {
	//跨域插件
	r.Use(Cors())
	rg := NewRegisterRoute(r, manager)
	rg.Register()
}
