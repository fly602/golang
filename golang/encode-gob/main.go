package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"runtime"
)

var (
	familyHashCacheFile = path.Join(os.Getenv("HOME"), ".cache", "deepin", "dde-daemon", "fonts", "family_hash")
)

type Family struct {
	Id   string
	Name string

	Styles []string

	Monospace bool
	Show      bool
}

type FamilyHashTable map[string]*Family

func main() {
	// 打开文件用于读取
	file, err := os.Open(familyHashCacheFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建解码器并将其关联到文件
	decoder := gob.NewDecoder(file)

	// 注册 Person 类型

	// 创建一个空的 Person 对象，用于存储反序列化后的数据
	table := make(FamilyHashTable)
	// 使用解码器将文件中的数据反序列化到 Person 对象
	err = decoder.Decode(&table)
	if err != nil {
		fmt.Println("Decode error:", err)
		return
	}

	// 输出反序列化后的 Person 对象
	for _, info := range table {
		fmt.Printf("Decoded family_hash:%+v\n", info)
	}

}

func printStackTrace() {
	// 获取调用栈信息
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)

	// 打印调用栈信息
	println("=== Call Stack ===")
	println(string(buf[:n]))
}
