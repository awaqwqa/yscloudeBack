package model

import (
	"crypto/rand"
	"fmt"
	"gorm.io/gorm"
)

const ShortUniqueHashLen = 16
const RandomHashLen = 32

type ShortUniqueHash string
type RedeemKey string
type RedeemStatus struct {
	Status string // newCreated, onGoing, finished
	Detail string
}

// 导入的表单

type BuildOption struct {
	Auth struct {
		Type  string `form:"type"`
		Token string `form:"token"`
	} `form:"auth"`
	// 任务名字
	TaskName string `form:"taskName"`
	// 操作的fbToken
	OptionalFBToken string `form:"optionalFBToken"`
	// 建筑地址
	StructureName string `form:"structureName"`
	// 服务器code
	RentalServerCode string `form:"rentalServerCode"`
	// 服务器密码
	RentalServerPassword string `form:"rentalServerPassword"`
	// x
	PosX int `form:"posX"`
	// y
	PosY int `form:"posY"`
	// z
	PosZ int `form:"posZ"`
}

// 任务栏
type BuildTaskInfo struct {
	Time        string
	InstanceID  string
	StartArgs   []string
	BuildOption *BuildOption
	FileSize    int64
}
type ProductionGroup struct {
	Name                     string
	ProductionGroupShortHash ShortUniqueHash
	AgentShortHash           ShortUniqueHash
	StructureNames           []string
	Redeems                  map[RedeemKey]*RedeemStatus
}

// 每个建筑的信息
type Structure struct {
	gorm.Model
	//用于锁定id
	StructureUserId uint
	//建筑文件名字
	FileName string
	//建筑文件类型 BDX SCHEMATIC ....
	FileType string
	//userName + 建筑文件的hash值
	FileHash string `gorm:"primary_key"`
	// 文件组 名字
	FileGroup string
	//建筑文件的大小
	FileSize int64
	//更新日期
	UploadDate string
}

func NewShortUniqueHash() ShortUniqueHash {
	return ShortUniqueHash(GenerateRandomHexStr(ShortUniqueHashLen))
}

// 用于生成指定长度的随机十六进制字符串
func GenerateRandomHexStr(len int) string {
	b := make([]byte, len/2+1)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:len]
}

func NewProductionGroup(syncDataFn func()) *ProductionGroup {
	return &ProductionGroup{
		StructureNames: []string{},
		Redeems:        map[RedeemKey]*RedeemStatus{},
	}
}
