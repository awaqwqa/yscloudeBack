package model

import (
	"crypto/rand"
	"fmt"
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
	TaskName             string `form:"taskName"`
	OptionalFBToken      string `form:"optionalFBToken"`
	StructureName        string `form:"structureName"`
	RentalServerCode     string `form:"rentalServerCode"`
	RentalServerPassword string `form:"rentalServerPassword"`
	PosX                 int    `form:"posX"`
	PosY                 int    `form:"posY"`
	PosZ                 int    `form:"posZ"`
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
	//建筑文件名字
	FileName string
	//建筑文件类型 BDX SCHEMATIC ....
	FileType string
	//建筑文件的hash值
	FileHash string
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
