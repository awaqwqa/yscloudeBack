package model

import (
	"crypto/rand"
	"fmt"
	"sync"
)

const ShortUniqueHashLen = 16
const RandomHashLen = 32

type ShortUniqueHash string
type RedeemKey string
type RedeemStatus struct {
	Status string // newCreated, onGoing, finished
	Detail string
}

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

type User struct {
	ID                int    `gorm:"primaryKey"`
	UserName          string `gorm:"ot null;unique;size:255"`
	Password          string `gorm:"not null"`
	Mobile            string `gorm:"unique;not null;"`
	mu                sync.Mutex
	FBToken           string
	Structures        map[string]*Structure
	Token             string
	ProductionGroups  map[ShortUniqueHash]*ProductionGroup
	ConsumedStructure int64
	RelatedInstances  map[string]*BuildTaskInfo
	ShortHash         ShortUniqueHash
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

func NewUser() *User {
	return &User{}
}
func (u *User) Get() {

}
func (u *User) UploadNewStructure(fileName string, fileData []byte) (err error) {

}
func (u *User) SetFBToken(fbToken string) {

}

// 获得fbtoken
func (a *User) GetFBToken() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.FBToken
}

// 取得所有的建筑信息备份
func (a *User) GetAllStructureInfoCopy() (structures map[string]Structure, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	structures = make(map[string]Structure)
	for k, v := range a.Structures {
		structures[k] = *v
	}
	return
}

// 添加一个导入组
func (a *User) AddProductionGroup(name string) (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if name == "" {
		return fmt.Errorf("production group name cannot be empty")
	}
	// check if name exists
	//查看是否存在
	for _, pg := range a.ProductionGroups {
		if pg.Name == name {
			return fmt.Errorf("production group name %v already exists", name)
		}
	}
	//创建一个新的导入组
	pg := NewProductionGroup(a.syncDataFn)
	pg.Name = name
	pg.AgentShortHash = a.ShortHash
	// allocate short hash for new production group
	for {
		shortHash := NewShortUniqueHash()
		if _, ok := a.ProductionGroups[shortHash]; !ok {
			pg.ProductionGroupShortHash = shortHash
			break
		}
	}
	a.ProductionGroups[ShortUniqueHash(pg.ProductionGroupShortHash)] = pg
	a.syncDataFn()
	return nil
}
