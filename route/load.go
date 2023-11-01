package route

import "yscloudeBack/source/app/controller"

func (rg *RegisterRoute) RegisterLoadRoute() error {
	rg.RegisterEngine.POST(joinRouterOnBasePath(getLoadUrl()), controller.LoadHandler(rg.Db))
	return nil
}
