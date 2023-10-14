package route

import (
	"github.com/gin-gonic/gin"
)

func RegisterLogRoute(r *gin.Engine) error {
	r.Group(getLogUrl())
	return nil
}
