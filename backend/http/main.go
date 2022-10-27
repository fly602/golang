package main

import "net/http"

func main() {
	// 绑定一个handler
	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir("./output"))))
	// 监听服务
	http.ListenAndServe(":8000", nil)
}
