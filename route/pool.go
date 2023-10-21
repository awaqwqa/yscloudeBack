package route

import "github.com/gin-gonic/gin"

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
func (rg *RegisterRoute) Register() {
	// TODO : connect the routers based on base_path to the handlefunction
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
func InitRoute(r *gin.Engine) {
	rg := NewRegisterRoute(r)
	rg.Register()

}
