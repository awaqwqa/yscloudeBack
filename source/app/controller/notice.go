package controller

import (
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func (cm *ControllerMannager) UpdateNotice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var form struct {
			Value string `json:"value"`
		}
		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		err = cm.dbManager.UpdateAnnouncement(1, form.Value)
		if err != nil {
			utils.Error(err.Error())
			model.BackErrorByString(ctx, err.Error())
			return
		}
		model.BackSuccess(ctx, nil)
	}
}
func (cm *ControllerMannager) GetNoticeValue() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		manager := cm.dbManager
		value, err := manager.GetAnnouncementByID(1)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		model.BackSuccess(ctx, value)
	}

}
