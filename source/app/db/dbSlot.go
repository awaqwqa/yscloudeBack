package db

import (
	"yscloudeBack/source/app/model"
)

func (dm *DbManager) CreateSlot(slot model.Slot) error {
	return dm.dbEngine.Create(&slot).Error
}

// 更新用户的年龄
func (dm *DbManager) UpdateSlotValue(slotStructId uint, value int) error {
	slot := model.Slot{}
	db := dm.dbEngine
	err := db.Model(&model.Slot{}).Where("id = ?", slotStructId).First(&slot).Error
	if err != nil {
		return err
	}

	slot.Value = value
	err = db.Model(&slot).Updates(slot).Error
	if err != nil {
		return err
	}

	return nil
}
func (dm *DbManager) GetAllSlots() ([]model.Slot, error) {
	var slots []model.Slot
	err := dm.dbEngine.Find(&slots).Error
	return slots, err
}

func (dm *DbManager) DeleteSlot(value int) error {
	return dm.dbEngine.Where("value = ?", value).Delete(&model.Slot{}).Error
}

// 根据slot id 找到对应slot
func (dm *DbManager) GetSlotById(id int) (model.Slot, error) {
	var slot model.Slot
	result := dm.dbEngine.Where("slot_id = ?", id).First(&slot)
	return slot, result.Error
}
