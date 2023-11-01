package controller

import (
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/model"
)

func LoadHandler(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lf *model.LoadForm
		model.BindStruct(ctx, lf)
		model.BackError(ctx, model.CodeUnknowError)
		return
		//TODO: 一些导入的处理
	}
}
