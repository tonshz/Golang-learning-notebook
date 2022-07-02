package global

import (
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var RootDir string

// 该类型的 Do 方法中的代码保证只会执行一次。
var once = new(sync.Once)

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

// inferRootDir 推断项目根目录
func inferRootDir() {
	// 通过 os.Getwd() 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 递归调用 infer
	// 判断目录 d 下面是否存在 template 目录
	var infer func(d string) string
	infer = func(d string) string {
		// 这里要确保项目根目录下存在 template 目录
		if exists(d + "/template") {
			return d
		}

		//如果 d 中不存在，则在其上级目录递归查找
		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
