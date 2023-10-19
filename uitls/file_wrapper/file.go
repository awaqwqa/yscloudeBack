package file_wrapper

import (
	"github.com/gin-gonic/gin"
	"sync"
)

func FileHandler(r *gin.Context) {

}

type File struct {
}

func NewFiler(rootDir string) *Filer {
	return &Filer{
		rootDir: rootDir,
	}
}

// 文件上传下载的函数
type Filer struct {
	rootDir string
	mu      sync.Mutex
	fileMap map[string]*File
}

func (f *Filer) InitFiler() {

}
func (f *Filer) NewFile(fileName string, data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()

}
func (f *Filer) DelFile(fileName string) (bool, error) {
	return true, nil
}
