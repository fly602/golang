package fileserver

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
)

func md5sum(reader io.Reader) string {
	// 创建MD5哈希对象
	hash := md5.New()
	// 从文件中读取数据并计算哈希值
	if _, err := io.Copy(hash, reader); err != nil {
		log.Fatal(err)
	}
	// 计算哈希值的字节表示
	hashBytes := hash.Sum(nil)

	// 将字节表示转换为十六进制字符串
	return fmt.Sprintf("%x", hashBytes)
}
