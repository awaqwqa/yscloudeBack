package model

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
)

// 这里需要注意的是Password是以md5加密储存的方式
type User struct {
	gorm.Model
	UserName string `gorm:"not null;unique;size:255"`
	//密码 md5(Password)
	Password          string     `gorm:"not null"`
	Mobile            string     `gorm:"not null;unique"`
	Key               string     `gorm:"not null;unique"`
	QQNumber          int        `gorm:"not null;unique"`
	mu                sync.Mutex `gorm:"-"`
	UserKeys          []Key      `gorm:"foreignKey:UserID"`
	FBToken           string
	Structures        map[string]*Structure `gorm:"-"`
	Token             string
	ProductionGroups  map[ShortUniqueHash]*ProductionGroup `gorm:"-"`
	ConsumedStructure int64                                `gorm:"-"`
	RelatedInstances  map[string]*BuildTaskInfo            `gorm:"-"`
	ShortHash         ShortUniqueHash                      `gorm:"-"`
}

func NewUser(name string, pwd string, key string) *User {
	return &User{
		Mobile:   key,
		UserName: name,
		Password: pwd,
		Key:      key,
		mu:       sync.Mutex{},
		UserKeys: []Key{},
	}
}
func (u *User) Get() {

}

func (u *User) GetLoadKeys() []Key {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.UserKeys
}
func (u *User) CheckLoadKey(key string) bool {
	for _, v := range u.UserKeys {
		fmt.Println(v.Value)
		if v.Value == key {
			return true
		}
	}
	return false
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
/*
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

*/
