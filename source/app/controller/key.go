package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"yscloudeBack/source/app/model"
)

func RegisterKey(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var k *model.Key
		if err := ctx.ShouldBindJSON(&k); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if ok {
				BackError(ctx, CodeUnknowError)
				return
			}
			fmt.Println("err:", err)
			BackError(ctx, CodeUnknowError)
			return
		}
		err := db.AddKey(k.Value)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": CodeSuccess,
			"msg":  CodeSuccess.Msg(),
		})
	}
}

// 获取keys
func GetKey(db *model.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keys, err := db.GetAllKeys()
		if err != nil {
			BackError(ctx, CodeUnknowError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": CodeSuccess,
			"msg":  CodeSuccess.Msg(),
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
			BackError(ctx, CodeInvalidKey)
			return
		}
		err := db.DeleteKey(form.DelKey)
		if err != nil {
			BackError(ctx, CodeUnknowError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": CodeSuccess,
			"msg":  CodeSuccess.Msg(),
		})
		return
	}

}
