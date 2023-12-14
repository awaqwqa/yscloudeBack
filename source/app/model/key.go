package model

import (
	"gorm.io/gorm"
	"yscloudeBack/source/app/utils"
)

const (
	USAGE_LOAD     = "导入"
	USAGE_REGISTER = "注册"
)

func NewKey(usage string, num int, fileGroupName string) (Key, error) {
	key, err := utils.GenerateRandomKey()
	if err != nil {
		return Key{}, err
	}
	return Key{
		Value:     key,
		Usage:     usage,
		Num:       num,
		FileGroup: fileGroupName,
		Status:    false,
	}, nil
}

// key的价格
type KeyPrice struct {
	ID    uint `gorm:"primaryKey"`
	Value int
}

// 作为密钥存储
type Key struct {
	gorm.Model
	// 用于锁定user
	UserID uint
	// 用途 比如导入 还是注册账号
	Usage string `json:"usage"`
	// 使用次数
	Num int `json:"num"`
	// 判断是否被用户获取
	Status bool `json:"isUsed"`
	//关联的file_group
	FileGroup string `json:"file_group"`
	Value     string `gorm:"size:32;unique;not null" json:"key" form:"key"`
}
