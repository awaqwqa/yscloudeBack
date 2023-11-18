package controller

import (
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func (cm *ControllerMannager) SetFbToken() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			FbToken string `json:"fbtoken"`
			//备注信息
			ReMark string `json:"remark"`
		}
		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		if manager.CheckFbToken(form.FbToken) {
			model.BackErrorByString(ctx, "this fbtoken is already exit")
			return
		}
		fbToken := model.NewFbToken(form.FbToken, form.ReMark)
		err = manager.AddFbToken(fbToken)
		if err != nil {
			utils.Error(err.Error())
			model.BackErrorByString(ctx, "add fbtoken false")
			return
		}
		model.BackSuccess(ctx, nil)
	}
}
func (cm *ControllerMannager) GetFbTokens() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		tokens, err := manager.GetFbTokens()
		if err != nil {
			model.BackErrorByString(ctx, "cant get fbtokens")
			return
		}

		model.BackSuccess(ctx, tokens)
	}
}
func (cm *ControllerMannager) DelFbTokens() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			DelToken string `json:"del_fbtoken"`
		}

		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		err = manager.DeleteFbToken(form.DelToken)
		if err != nil {
			model.BackErrorByString(ctx, "del fbtoken false")
			return
		}
		model.BackSuccess(ctx, nil)
	}
}
