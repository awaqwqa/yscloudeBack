package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindStruct(ctx *gin.Context, st any) (MyCode, error) {
	if err := ctx.ShouldBindJSON(&st); err != nil {
		_, ok := err.(validator.ValidationErrors)
		if ok {
			return CodeUnknowError, fmt.Errorf("Validation")
		}
		fmt.Println("err:", err)
		return CodeUnknowError, fmt.Errorf("unkonw")
	}
	return CodeSuccess, nil
}
