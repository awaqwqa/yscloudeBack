package route

const (
	BASE_PATH = "/v1/ysback"
	LOGPATH   = "auth"
	REGISTER  = "register"
	LOGIN     = "login"
	LOGOUT    = "logout"
	LOADPATH  = "load"
	FILEPATH  = "file"
)

func getBaseUrl() string {
	return BASE_PATH
}
func getLogUrl() string {
	return LOGPATH
}
func getRegisterUrl() string {
	return REGISTER
}
func getLoginUrl() string {
	return LOGIN
}
func getLogoutUrl() string {
	return LOGOUT
}
func getLoadUrl() string {
	return LOADPATH
}
func getFileUrl() string {
	return FILEPATH
}
func joinRouter(args ...string) (newRoute string) {
	for _, v := range args {
		newRoute += "/" + v
	}
	return
}
func joinRouterOnBasePath(args ...string) (newRouter string) {
	newRouter = getBaseUrl()
	for _, v := range args {
		newRouter += "/" + v
	}
	return
}
func getLogUrlOnBase() string {
	return joinRouterOnBasePath(getLogUrl())
}
