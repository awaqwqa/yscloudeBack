package model

import (
	"gorm.io/gorm"
	"sync"
	"yscloudeBack/source/app/utils"
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
	UserKeys          []UserKey
	FBToken           string
	Structures        map[string]*Structure `gorm:"-"`
	Token             string
	ProductionGroups  map[ShortUniqueHash]*ProductionGroup `gorm:"-"`
	ConsumedStructure int64                                `gorm:"-"`
	RelatedInstances  map[string]*BuildTaskInfo            `gorm:"-"`
	ShortHash         ShortUniqueHash                      `gorm:"-"`
}
type UserKey struct {
	gorm.Model
	UserID    uint
	Value     string `json:"key"`
	Usage     string `json:"usage"`
	Num       int    `json:"num"`
	Status    bool   `json:"isUsed"`
	FileGroup string `json:"file_group"`
}

func NewUser(name string, pwd string, key string) *User {
	return &User{
		Mobile:   key,
		UserName: name,
		Password: pwd,
		Key:      key,
		mu:       sync.Mutex{},
		UserKeys: []UserKey{},
	}
}
func (u *User) Get() {

}
func (u *User) GetKeys() []UserKey {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Update()
	return u.UserKeys
}
func (u *User) DelKeys(key string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	for k, v := range u.UserKeys {
		if v.Value != key {
			continue
		}

		newSlice, err := utils.RemoveIndex(u.UserKeys, k)
		if err != nil {
			return false
		}
		u.UserKeys = newSlice
		u.Upgrade()
		return true
	}
	return false
}

// 从数据库获取最新数据
func (u *User) Update() {

}

// 将内容更新的数据库
func (u *User) Upgrade() {

}

/*
func (u *User) UploadNewStructure(fileName string, fileData []byte) (err error) {

}

*/

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
