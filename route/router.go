package route

const (
	BASE_PATH = "/v1/ysback"
	LOGPATH   = "/auth"
)

func getBaseUrl() string {
	return BASE_PATH
}
func getLogUrl() string {
	return LOGPATH
}
func joinRouter(args ...string) (newRoute string) {
	for _, v := range args {
		newRoute += v
	}
	return
}
func joinRouterOnBasePath(args ...string) (newRouter string) {
	newRouter = getBaseUrl()
	for _, v := range args {
		newRouter += v
	}
	return
}
func getLogUrlOnBase() string {
	return joinRouterOnBasePath(getLogUrl())
}
