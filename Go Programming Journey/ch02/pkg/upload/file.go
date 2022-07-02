package upload

import (
	"demo/ch02/global"
	"demo/ch02/pkg/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 定义 FileType 为 int 的别名
type FileType int

// 使用 FileType 作为类别表示的基础类型
const TypeImage FileType = iota + 1 // TypeImage = 1

// 返回经过加密处理后的文件名
func GetFileName(name string) string {
	ext := GetFileExt(name)
	// 返回没有后缀的文件名
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

func GetFileExt(name string) string {
	// 调用 path.Ext() 进行循环查找
	return path.Ext(name)
}

func GetSavePath() string {
	// 返回配置中的文件保存目录
	return global.AppSetting.UploadSavePath
}

// 检测路径是否存在
func CheckSavePath(dst string) bool {
	// 利用 os.Stat() 方法所返回的 error 值与系统所定义的 oserror.ErrNotExist 是否相等
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

// 检测文件后缀是否满足设置条件
func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	// 同一转换为大写进行匹配
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		// 与配置文件中设置的允许的文件后缀名进行比较
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}
	}
	return false
}

// 检测最大大小是否超出最大限制
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}
	return false
}

// 检测文件权限是否足够
func CheckPermission(dst string) bool {
	// 与 CheckSavePath() 类似，与 oserror.ErrPermission 判断是否相等
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

// 创建上传文件时所使用的保存目录
func CreateSavePath(dst string, perm os.FileMode) error {
	// os.MkdirAll 会根据传入的 os.FileMode 权限位递归创所需的所有目录结构
	// 若目录已存在则不会进行任何操作，直接返回 nil
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

// 保存所上传的文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	// 通过 file.Open 打开源地址的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// 通过 os.Create 创建目标地址的文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	// 结合 io.Copy 实现两者之间的文件内容拷贝
	_, err = io.Copy(out, src)
	return err
}
