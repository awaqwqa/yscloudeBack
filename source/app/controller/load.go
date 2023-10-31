package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"yscloudeBack/source/app/model"
)

func LoadHandler(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lf *model.LoadForm
		if err := ctx.ShouldBindJSON(&lf); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if ok {
				BackError(ctx, CodeUnknowError)
				return
			}
			fmt.Println("err:", err)
			BackError(ctx, CodeUnknowError)
			return
		}
		BackError(ctx, CodeUnknowError)
		return
		//TODO: 一些导入的处理
	}
}
