package model

import "gorm.io/gorm"

type FbToken struct {
	gorm.Model
	Value string
}

func NewFbToken(value string) *FbToken {
	return &FbToken{
		Value: value,
	}
}
