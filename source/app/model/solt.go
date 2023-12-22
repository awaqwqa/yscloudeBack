package model

import (
	"gorm.io/gorm"
	"yscloudeBack/source/app/utils"
)

const (
	DISPOSABLE = "disposable"
	PERMANENT  = "permanent"
)

func NewSlot(userID uint, soltType string, Value int) Slot {
	return Slot{
		UserID:   userID,
		SlotType: soltType,
		Value:    Value,
		SlotId:   int(utils.GenerateUniqueIntID()),
	}
}

type Slot struct {
	gorm.Model
	UserID   uint
	SlotType string `json:"slot_type"`
	Value    int    `json:"slot_value"`
	Time     int64  `json:"slot_time"`
	// 查询值
	SlotId int `json:"slot_id"`
}
