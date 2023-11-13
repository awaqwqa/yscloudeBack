package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"yscloudeBack/source/app/utils"
)

func BindStruct(ctx *gin.Context, st any) (MyCode, error) {
	if err := ctx.ShouldBindJSON(&st); err != nil {
		_, ok := err.(validator.ValidationErrors)
		if ok {
			return CodeUnknowError, fmt.Errorf("Validation")
		}
		utils.Error(err.Error())
		return CodeBindFalse, fmt.Errorf("unkonw")
	}
	return CodeBindFalse, nil
}
