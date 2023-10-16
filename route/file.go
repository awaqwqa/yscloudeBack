package route

import "yscloudeBack/source/app/controller"

func (rg *RegisterRoute) RegisterFileRoute() error {
	rg.RegisterEngine.POST(joinRouterOnBasePath(getFileUrl()), controller.FileHandler)
	return nil
}
