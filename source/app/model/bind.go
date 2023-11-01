package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindStruct(ctx *gin.Context, st any) {
	if err := ctx.ShouldBindJSON(&st); err != nil {
		_, ok := err.(validator.ValidationErrors)
		if ok {
			BackError(ctx, CodeUnknowError)
			return
		}
		fmt.Println("err:", err)
		BackError(ctx, CodeUnknowError)
		return
	}
}
