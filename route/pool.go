package route

import "github.com/gin-gonic/gin"

type handleFunction func(ctx *gin.Context)

func getHandleAndRouter() map[string]handleFunction {
	return map[string]handleFunction{}
}
func NewRegisterRoute(rg *gin.RouterGroup) *RegisterRoute {
	return &RegisterRoute{
		RegisterGroup: rg,
	}
}

type RegisterRoute struct {
	RegisterGroup *gin.RouterGroup
}

// Initialization is performed (执行) to connect the router to the handle function.
func (rg *RegisterRoute) Register(DicRouterPool map[string]handleFunction) {
	// TODO : connect the routers based on base_path to the handlefunction
}

// router wouldnt be imported .The router package is used to initialize the router similar a controller
func InitRoute(r *gin.Engine) {
	RouterGroup := r.Group(getBaseUrl())
	rg := NewRegisterRoute(RouterGroup)
	rg.Register(getHandleAndRouter())
}
