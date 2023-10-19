package route

import (
	"yscloudeBack/uitls/file_wrapper"
)

func (rg *RegisterRoute) RegisterFileRoute() error {
	rg.RegisterEngine.POST(joinRouterOnBasePath(getFileUrl()), file_wrapper.FileHandler)
	return nil
}
