package archiveManager

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"sync"
)

type InstanceDetail struct {
	InstanceID   string   `json:"instance_id"`
	Name         string   `json:"name"`
	Cmd          string   `json:"cmd"`
	Args         []string `json:"args"`
	Status       string   `json:"status"`
	StatusDetail string   `json:"status_detail"`
}

// 建筑管理器
// 存储有archiveMap类似这样:"instance.006570ae639040fb9569a1d6ac7dbe1a.detail": "archive/aW5zdGFuY2UuMDA2NTcwYWU2MzkwNDBmYjk1NjlhMWQ2YWM3ZGJlMWEuZGV0YWls"
// 还含有存储根目录
// 存储方法
// 锁
type ArchiveManager struct {
	//储存
	storage string
	//建筑对应dic
	archiveMap map[string]string
	syncDataFn func()
	//锁
	mu sync.Mutex
}

// 一个新的建筑管理器
func NewArchiveManager(storageDir string) *ArchiveManager {
	//初始化
	a := &ArchiveManager{storage: storageDir, mu: sync.Mutex{}}
	//创建一个./archive目录
	err := os.MkdirAll(storageDir, 0755)
	if err != nil {
		return nil
	}
	//获取archiveMap信息
	jsonFile := path.Join(storageDir, "archiveMap.json")
	//初始化表格
	archiveMap := map[string]string{}
	//安全读取指定目录下方的json文件并且写入结构体
	//所以数据库得有个archiveMap表格 存有key 与vaule两个值
	//json中类似这样"instance.006570ae639040fb9569a1d6ac7dbe1a.detail": "archive/aW5zdGFuY2UuMDA2NTcwYWU2MzkwNDBmYjk1NjlhMWQ2YWM3ZGJlMWEuZGV0YWls"
	if err := LoadJson(jsonFile, &archiveMap); err != nil {
		fmt.Println("load archiveMap.json err:", err)
	}
	//同步
	a.archiveMap = archiveMap
	//保存json文件
	a.syncDataFn = func() {
		SaveJson(jsonFile, a.archiveMap)
	}
	return a
}

// 创建一个新的建筑
func (a *ArchiveManager) ArchiveNew(name string, data []byte) {
	fName := path.Join(a.storage, base64.StdEncoding.EncodeToString([]byte(name)))
	os.WriteFile(fName, data, 0755)
	a.mu.Lock()
	a.syncDataFn()
	defer a.mu.Unlock()
	a.archiveMap[name] = fName
}

func (a *ArchiveManager) GetArchive(name string) (data []byte, err error) {
	a.mu.Lock()
	// 从archiveMap中获取文件名 instance.%v.detail
	fname, found := a.archiveMap[name]
	a.mu.Unlock()
	if !found {
		return nil, fmt.Errorf("not found")
	} else {
		file, err := os.ReadFile(fname)
		return file, err
	}
}
