package filer

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"yscloudeBack/source/app/utils"
)

// 判断路径是绝对地址还是相对地址
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// 检查名字是否符合路径命名标准
func IsNameValid(name string) bool {
	// 使用正则表达式进行匹配
	pattern := `^[\w\-. ]+$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(name)
}

func GetAbsolutePath() (string, error) {
	absPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// 获取当前文件的相对地址（相对于工作目录）
func GetRelativePath() (string, error) {
	absPath, err := GetAbsolutePath()
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(".", absPath)
	if err != nil {
		return "", err
	}

	return relPath, nil
}

// 向指定路径写入文件
func WriteToFile(filePath string, content []byte) error {
	err := ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

// 判断指定路径是否合法
func IsPathValid(name string) bool {
	// 检查名字是否为空
	if name == "" {
		return false
	}

	// 使用正则表达式检查名字的格式
	// 这里使用的是 Windows 文件/文件夹命名规范的正则表达式
	regex := regexp.MustCompile(`^[^\x00-\x1F<>:"/\\|?*]+$`)
	return regex.MatchString(name)
}

// 获取指定文件夹下的信息
func GetFolderInfo(folderPath string) ([]os.FileInfo, error) {
	utils.Info("get folder info : folderPath >> %v", folderPath)
	fileInfo, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

// 判断路径是文件夹还是文件
func GetPathType(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		return "folder" // 文件夹
	}

	return "file" // 文件
}

// 判断路径是否存在
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false // 路径不存在
		}
	}
	return true
}

// 读取文件内容
func ReadFileContent(filePath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// 删除文件
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
func DeleteDir(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}

	return nil
}

// 创建文件夹
func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
