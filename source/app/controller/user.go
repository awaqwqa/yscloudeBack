package controller

import "C"
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func checkName(name string) (bool, error) {
	return false, nil
}

// Register 用户注册
func Register(ctx *gin.Context) {
	name := ctx.PostForm("name")
	key := ctx.PostForm("key")
	passwd := ctx.PostForm("passwd")
	if ok, err := checkName(name); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": err.Error()})
	}
}
func Login(ctx *gin.Context) {

}
