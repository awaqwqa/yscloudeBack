package model

import "gorm.io/gorm"

type Key struct {
	gorm.Model
	Value string `gorm:"size:32;unique;not null" json:"key" form:"key"`
}
