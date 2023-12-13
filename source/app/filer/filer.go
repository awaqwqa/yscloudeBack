package filer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"yscloudeBack/source/app/utils"
)

type UserDirName string
type taskData []byte
type GroupName string

type Task struct {
	// 这里是操作的文件名字
	FileName  string
	FileGroup string
}
type Filer struct {
	rootPath string
	//这里存放的是用户的地址
	taskPool map[UserDirName][]Task
	sync.Mutex
}

func NewFiler(rootPath string) *Filer {
	return &Filer{
		rootPath: rootPath,
		taskPool: make(map[UserDirName][]Task),
	}
}
func (f *Filer) GetUserPath(userName string) string {
	return filepath.Join(f.rootPath, userName)
}

func CheckUserPath(path string) error {
	if !IsPathExists(path) {
		return fmt.Errorf("user dir is not exit")
	}
	return nil
}

func CheckFileGroupPath(user_group_name string) error {
	if !IsPathExists(user_group_name) || GetPathType(user_group_name) != "folder" {
		return fmt.Errorf("file_group dir is not exit")
	}
	return nil
}
func CheckFilePath(file_path string) error {
	if !IsPathExists(file_path) || GetPathType(file_path) != "file" {
		return fmt.Errorf("file path is not exit")
	}
	return nil
}

// 创建一个存储任务
func (f *Filer) NewTask(userName string, groupName string, fileName string, data []byte) (err error) {
	var cancelFunc func()
	// 用于清理任务
	cancelFunc = func() {
		f.Lock()
		defer f.Unlock()
		// 还需要做一些意外处理
		utils.Warning("cancelFunc start")
		if _, ok := f.taskPool[UserDirName(userName)]; !ok {
			return
		}
		for k, v := range f.taskPool[UserDirName(userName)] {
			if v.FileName == fileName {
				slice, err2 := utils.RemoveIndex(f.taskPool[UserDirName(userName)], k)
				if err2 != nil {
					utils.Error(err2.Error())
					return
				}
				f.taskPool[UserDirName(userName)] = slice
			}
		}
	}
	{
		f.Lock()
		if _, ok := f.taskPool[UserDirName(userName)]; !ok {
			f.taskPool[UserDirName(userName)] = []Task{}
		}
		for _, v := range f.taskPool[UserDirName(userName)] {
			if v.FileName == fileName {
				cancelFunc()
				return fmt.Errorf("this task allready exit,dont create again")
			}
		}

		f.taskPool[UserDirName(userName)] = append(f.taskPool[UserDirName(userName)], Task{FileName: fileName, FileGroup: groupName})
		f.Unlock()
	}
	user_path := path.Join(f.rootPath, userName)
	err = CheckUserPath(user_path)
	if err != nil {
		cancelFunc()
		return err
	}
	file_group_path := path.Join(user_path, groupName)
	err = CheckFileGroupPath(file_group_path)
	if err != nil {
		cancelFunc()
		return err
	}
	file_path := path.Join(file_group_path, fileName)

	if !IsPathValid(fileName) {
		cancelFunc()
		return fmt.Errorf("create a new task but this path is not allowed")
	}
	if !IsPathValid(fileName) {
		cancelFunc()
		return fmt.Errorf("this file path is not allowed")
	}

	// 开始执行写入任务
	go func() {
		defer cancelFunc()
		err := WriteToFile(file_path, data)
		if err != nil {
			utils.Error(err.Error())
			return
		}
	}()
	return nil
}

// 获取指定玩家的文件组中文件数据
func (f *Filer) GetFileGroupFileData(name string, fileGroupName string, fileName string) ([]byte, error) {
	user_path := path.Join(f.rootPath, name)
	err := CheckUserPath(user_path)
	if err != nil {
		return nil, err
	}
	user_group_name := path.Join(user_path, fileGroupName)
	err = CheckFileGroupPath(user_group_name)
	if err != nil {
		return nil, err
	}
	file_path := path.Join(user_group_name, fileName)
	if !IsPathExists(file_path) || GetPathType(file_path) != "file" {
		return nil, fmt.Errorf("file path is not exit")
	}
	content, err := ReadFileContent(file_path)
	if err != nil {
		return nil, err
	}
	return content, nil

}

// 获取指定文件组下的文件头信息
func (f *Filer) GetFileGroupHeadData(name string, fileGroupName string) ([]os.FileInfo, error) {
	user_path := path.Join(f.rootPath, name)
	if !IsPathExists(user_path) {
		return nil, fmt.Errorf("user dir is not exit")
	}
	user_group_name := path.Join(user_path, fileGroupName)
	if !IsPathExists(user_group_name) || GetPathType(user_group_name) != "folder" {
		return nil, fmt.Errorf("file_group dir is not exit")
	}
	infos, err := GetFolderInfo(user_group_name)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// 根据名字获取 文件组名字
func (f *Filer) GetFileGroupsName(name string) ([]GroupName, error) {
	file_path := path.Join(f.rootPath, name)
	utils.Warning("get file groups : file_path >> ", file_path)
	if !IsPathExists(file_path) {
		utils.Error("user dir is not exit , allreadly created yet")
		err := CreateFolder(file_path)
		if err != nil {
			return nil, err
		}
	}
	infos, err := GetFolderInfo(file_path)
	if err != nil {
		return nil, err
	}
	group_names := []GroupName{}
	for _, v := range infos {
		if !v.IsDir() {
			continue
		}
		group_names = append(group_names, GroupName(v.Name()))
	}
	return group_names, nil
}
func (f *Filer) CreateFileGroupDir(name string, fileGroupName string) error {
	user_path := filepath.Join(f.rootPath, name)
	utils.Info("create file group : user_path >> %v", user_path)
	err := CheckUserPath(user_path)
	if err != nil {
		utils.Error(err.Error())
		err = nil
		if !IsPathValid(name) {
			return fmt.Errorf("system want to create a user dir , but this user_path is not allowed")
		}
		err := f.CreateUserDir(name)
		if err != nil {
			return err
		}
	}
	file_group_path := filepath.Join(user_path, fileGroupName)
	if !IsPathValid(fileGroupName) {
		utils.Error(file_group_path)
		return fmt.Errorf("system want to create a file_group_path, bur this path is not allowed")
	}
	err = CreateFolder(file_group_path)
	if err != nil {
		utils.Error(err.Error())
		return fmt.Errorf("create file_group dir false")
	}
	return nil
}

// 创建用户文件夹
func (f *Filer) CreateUserDir(name string) error {
	user_path := path.Join(f.rootPath, name)
	utils.Warning("check user_path :", user_path)
	if IsPathExists(user_path) {
		return fmt.Errorf("already exit this user dir")
	}
	if !IsPathValid(name) {
		return fmt.Errorf("this path not be allowed")
	}
	err := CreateFolder(user_path)
	if err != nil {
		utils.Error(err.Error())
		return fmt.Errorf("create user_path false")
	}
	return nil
}
func (f *Filer) DelectUserDir(name string) error {
	user_path := path.Join(f.rootPath, name)
	if !IsPathExists(user_path) || GetPathType(user_path) != "folder" {
		return fmt.Errorf("this user dir is not exit")
	}
	err := DeleteDir(user_path)
	if err != nil {
		return err
	}
	return nil
}
func (f *Filer) DelectFileGroup(name string, fileGroupName string) error {
	user_path := filepath.Join(f.rootPath, name)
	err := CheckUserPath(user_path)
	if err != nil {
		return err
	}
	file_group_path := filepath.Join(user_path, fileGroupName)
	err = CheckFileGroupPath(file_group_path)
	if err != nil {
		return err
	}
	for _, v := range f.taskPool[UserDirName(name)] {
		if v.FileGroup == file_group_path {
			return fmt.Errorf("%v file group is working,you cant delect it", fileGroupName)
		}
	}
	err = DeleteDir(file_group_path)
	if err != nil {
		utils.Error(err.Error())
		return fmt.Errorf("errors: cant delect file group.reason: unkown")
	}
	return nil
}
