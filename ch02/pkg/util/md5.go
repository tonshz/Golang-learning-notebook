package util

import (
	"crypto/md5"
	"encoding/hex"
)

// 将上传后的文件名进行格式化，将文件名 MD5 编码后在进行写入
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}
