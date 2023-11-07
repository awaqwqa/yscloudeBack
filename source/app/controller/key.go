package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func RegisterKey(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, err := utils.GenerateRandomKey()
		if err != nil {
			model.BackError(ctx, model.CodeGetKeyFalse)
			return
		}
		err = db.AddKey(key)
		if err != nil {
			model.BackError(ctx, model.CodeInvalidKey)
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
		key_values := []string{}
		for _, v := range keys {
			key_values = append(key_values, v.Value)
		}
		model.BackSuccess(ctx, key_values)
		return
	}
}
func DelKey(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var form struct {
			DelKey string `form:"del_key" binding:"required" json:"del_key"`
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
