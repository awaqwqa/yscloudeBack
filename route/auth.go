package route

import (
	"yscloudeBack/source/app/controller"
)

func (rg *RegisterRoute) RegisterLogRoute() error {
	logGroup := rg.RegisterEngine.Group(joinRouterOnBasePath(getLogUrl()))
	{
		// Register
		logGroup.POST(getRegisterUrl(), controller.Register(rg.Db))
		//logGroup.POST(getLoginUrl(), controller.Login)
		//logGroup.POST(getLogoutUrl())
	}
	return nil
}
