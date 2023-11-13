package controller

import (
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/model"
)

func LoadHandler(db *db.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lf *model.LoadForm
		code, err := model.BindStruct(ctx, &lf)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		model.BackError(ctx, model.CodeUnknowError)
		return
		//TODO: 一些导入的处理
	}
}
