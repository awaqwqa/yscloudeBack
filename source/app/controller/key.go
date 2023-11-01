package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yscloudeBack/source/app/model"
)

func RegisterKey(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var k *model.Key
		code, err2 := model.BindStruct(ctx, k)
		if err2 != nil {
			model.BackError(ctx, code)
			return
		}
		err := db.AddKey(k.Value)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": model.CodeSuccess,
			"msg":  model.CodeSuccess.Msg(),
		})
	}
}

// 获取keys
func GetKey(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keys, err := db.GetAllKeys()
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": model.CodeSuccess,
			"msg":  model.CodeSuccess.Msg(),
			"keys": keys,
		})
		return
	}
}
func DelKey(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var form struct {
			DelKey string `form:"del_key" binding:"required"`
		}

		if err := ctx.ShouldBind(&form); err != nil {
			model.BackError(ctx, model.CodeInvalidKey)
			return
		}
		err := db.DeleteKey(form.DelKey)
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": model.CodeSuccess,
			"msg":  model.CodeSuccess.Msg(),
		})
		return
	}

}
