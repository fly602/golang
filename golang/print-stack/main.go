package main

import (
	"fmt"
	"runtime"
)

func main() {
	// 打印调用堆栈
	printStackTrace()

	// 这里可以添加你的其他代码
}

func printStackTrace() {
	// 获取调用堆栈信息
	stack := make([]byte, 8192)
	length := runtime.Stack(stack, false)

	// 打印调用堆栈
	fmt.Printf("=== Call Stack ===\n%s\n", string(stack[:length]))
}
