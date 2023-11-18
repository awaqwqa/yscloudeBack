package model

import "gorm.io/gorm"

type FbToken struct {
	gorm.Model
	Value  string
	ReMark string
}

func NewFbToken(value string, remark string) *FbToken {
	return &FbToken{
		Value:  value,
		ReMark: remark,
	}
}
