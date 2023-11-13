package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/model"
)

// 注册key格式:
//
//	{
//		"usage":"导入",
//		"num":1,
//		"fileGroupName":""
//	}
func RegisterKey(db *db.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var form struct {
			Usage string `json:"usage"`
			Num   int    `json:"num"`
		}
		code, err2 := model.BindStruct(ctx, &form)
		if err2 != nil {
			model.BackError(ctx, code)
			return
		}
		if form.Usage != model.USAGE_LOAD && form.Usage != model.USAGE_REGISTER {
			model.BackSuccess(ctx, model.CodeCodeTypeFalse)
			return
		}
		if !(form.Num > 0) {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		key, err := model.NewKey(form.Usage, form.Num, "")
		if err != nil {
			model.BackError(ctx, model.CodeCreateKeyFalse)
			return
		}
		err = db.AddKey(key)
		if err != nil {
			fmt.Println(err)
			model.BackError(ctx, model.CodeInvalidKey)
			return
		}
		value, err := db.GetKeyByValue(key.Value)
		if err != nil {
			model.BackError(ctx, model.CodeGetKeyFalse)
			return
		}
		model.BackSuccess(ctx, value)
	}
}

// 获取keys
func GetKey(db *db.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keys, err := db.GetAllKeys()
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		//key_values := []string{}
		//for _, v := range keys {
		//	key_values = append(key_values, v.Value)
		//}
		model.BackSuccess(ctx, keys)
		return
	}
}
func DelKey(db *db.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var form struct {
			DelKey string `form:"del_key" binding:"required" json:"del_key"`
		}

		if err := ctx.ShouldBind(&form); err != nil {

			model.BackError(ctx, model.CodeInvalidKey)
			return
		}
		isFind, key := db.CheckKey(form.DelKey)
		if key != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		if !isFind {
			model.BackError(ctx, model.CodeGetKeyFalse)
			return
		}

		err := db.DeleteKey(form.DelKey)
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		model.BackSuccess(ctx, nil)
		return
	}

}
