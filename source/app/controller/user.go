package controller

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"yscloudeBack/source/app/model"
)

func checkName(name string) (bool, error) {
	return false, nil
}

// Register 用户注册
func Register(ctx *gin.Context) {
	var r *model.RegisterForm
	if err := ctx.ShouldBindJSON(&r); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			
		}
	}
	name := ctx.PostForm("name")
	key := ctx.PostForm("key")
	passwd := ctx.PostForm("passwd")
	if ok, err := checkName(name); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": err.Error()})
	}
}
func Login(ctx *gin.Context) {

}
