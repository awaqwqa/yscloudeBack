package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/middleware"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func (cm *ControllerMannager) AddSlots() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			Key string `json:"Key"`
		}
		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		name, err := middleware.GetContextName(ctx)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		user, err := manager.GetUserByUserName(name)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		//	TODO: 查询key类型
		key, err := manager.GetKeyByValue(form.Key)
		if err != nil {
			utils.Error(err.Error())
			model.BackErrorByString(ctx, fmt.Sprintf("get key type false"))
			return
		}
		if key.Usage != model.USAGE_PRSLOT && key.Usage != model.USAGE_DISSLOT {
			model.BackErrorByString(ctx, fmt.Sprintf("key usage is not used for slot"))
			return
		}
		slot := model.NewSlot(user.ID, key.Usage, 0)
		err = manager.DeleteKey(form.Key)
		if err != nil {
			model.BackErrorByString(ctx, fmt.Sprintf("delect key false"))
			return
		}
		err = manager.CreateSlot(slot)
		if err != nil {
			utils.Error(err.Error())
			model.BackErrorByString(ctx, err.Error())
			return
		}
		model.BackSuccess(ctx, fmt.Sprintf("new slot success"))
	}
}
func (cm *ControllerMannager) DelUserSlots() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			SlotId int `json:"slot_id"`
		}
		code, err2 := model.BindStruct(ctx, &form)
		if err2 != nil {
			model.BackError(ctx, code)
			return
		}
		user, err := cm.GetUserFromCtx(ctx)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		for _, v := range user.Slots {
			if v.SlotId == form.SlotId {
				err := manager.DeleteSlot(v.Value)
				if err != nil {
					model.BackErrorByString(ctx, fmt.Sprintf("delete slot false"))
					return
				}
				model.BackSuccess(ctx, fmt.Sprintf("success del %v slot", form.SlotId))
				return
			}
		}
		model.BackErrorByString(ctx, fmt.Sprintf("the slots id does not exist"))
	}
}
func (cm *ControllerMannager) GetUserSlots() gin.HandlerFunc {
	//manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		user, err := cm.GetUserFromCtx(ctx)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		//fmt.Println("user.slots", user.Slots)
		model.BackSuccess(ctx, user.Slots)
	}
}
func (cm *ControllerMannager) UpdateSlots() gin.HandlerFunc {
	manager := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			SlotValue int `json:"slot_value"`
			SlotId    int `json:"slot_id"`
		}

		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		if form.SlotValue == 0 || form.SlotId == 0 {
			model.BackErrorByString(ctx, fmt.Sprintf("update slot with Value and Id!!!!"))
			return
		}
		user, err := cm.GetUserFromCtx(ctx)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		// 获取user id
		isFind := false
		slot_index := 0
		func() {
			for k, v := range user.Slots {
				fmt.Println("slot_id:", v.SlotId)
				if v.SlotId == form.SlotId {

					isFind = true
					slot_index = k
				}
			}
		}()
		if !isFind {
			model.BackErrorByString(ctx, fmt.Sprintf("cant get %v id slot", form.SlotId))
			return
		}

		//if slot.UserID != userId {
		//	model.BackErrorByString(ctx, fmt.Sprintf("You do not have permission to use this slot"))
		//	return
		//}
		slot := user.Slots[slot_index]
		if slot.SlotType == model.DISPOSABLE {
			model.BackErrorByString(ctx, fmt.Sprintf("this slot is disposable ,cant be change"))
			return
		}
		slot.Value = form.SlotValue
		err = manager.UpdateSlotValue(slot.ID, slot.Value)
		if err != nil {
			utils.Error(err.Error())
			model.BackErrorByString(ctx, fmt.Sprintf("when update slot,insert Slot false"))
			return
		}
		model.BackSuccess(ctx, fmt.Sprintf("The value of %v slot has been changed to %v ", slot.SlotId, slot.Value))
	}
}
