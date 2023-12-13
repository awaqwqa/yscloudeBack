package model

import (
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

// 这里需要注意的是Password是以md5加密储存的方式
type User struct {
	gorm.Model
	UserName string `gorm:"not null;unique;size:255"`
	//密码 md5(Password)
	Password   string     `gorm:"not null"`
	Mobile     string     `gorm:"not null;unique"`
	Key        string     `gorm:"not null;unique"`
	QQNumber   int        `gorm:"not null;unique"`
	mu         sync.Mutex `gorm:"-"`
	UserKeys   []Key      `gorm:"foreignKey:UserID"`
	FBToken    string
	Structures []Structure `gorm:"foreignKey:StructureUserId"`
	// 余额值
	Balance           int `gorm:"balance"`
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

func (a *User) NewUserStructure(fileName string, fileData []byte, file_group string) (structure Structure, err error) {
	// overwrite if exists
	// check file type: BDX if .bdx file, Schem if .schem file, schematic if .schematic file
	// any other file type will be rejected
	fileType := "unknown"
	if path.Ext(fileName) == ".bdx" {
		fileType = "BDX"
	} else if path.Ext(fileName) == ".schem" {
		fileType = "Schem"
	} else if path.Ext(fileName) == ".schematic" {
		fileType = "schematic"
	} else if path.Ext(fileName) == ".mcworld" {
		fileType = "mcworld"
	} else {
		return structure, fmt.Errorf("file type not supported")
	}

	// save file
	filePath := a.GetLoadPath(fileName)
	if err = ioutil.WriteFile(filePath, fileData, 0755); err != nil {
		return
	}

	structure = Structure{
		FileName:   fileName,
		FileType:   fileType,
		FileHash:   GetFileHash(a.UserName, fileData),
		FileSize:   int64(len(fileData)),
		UploadDate: GetNowString(),
		FileGroup:  file_group,
	}
	return
}

// 取得所有的建筑信息备份
func (a *User) GetAllStructureInfoCopy() (structures []Structure, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	structures = []Structure{}
	for _, v := range a.Structures {
		structures = append(structures, v)
	}
	return
}

// 取得user的存文件的相对地址
func (a *User) GetDirPath() string {
	workDir, _ := os.Getwd()
	return path.Join(workDir, a.UserName)
}
func (a *User) GetLoadPath(structName string) string {
	return path.Join(a.GetDirPath(), structName)
}
func (a *User) ReadStructs() (structures []Structure, err error) {
	f, err := os.Open(a.GetDirPath())
	if err != nil {
		return
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			structures = append(structures, Structure{
				FileName: file.Name(),
				FileSize: file.Size(),
				FileType: filepath.Ext(file.Name()),
			})
		}
	}

	return structures, nil
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
func GetFileHash(name string, data []byte) string {
	return fmt.Sprintf("%x%x", name, md5.Sum(data))
}
func GetNowString() string {
	nt := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", nt.Year(), nt.Month(), nt.Day())
}
